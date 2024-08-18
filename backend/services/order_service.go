package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type OrderService struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

type OrderStatistics struct {
	Month  string `json:"month"`
	Count  int    `json:"count"`
	Status string `json:"status,omitempty"`
}

func NewOrderService(collection *mongo.Collection) *OrderService {
	return &OrderService{
		collection:  collection,
		redisClient: utils.RedisClient,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) (*mongo.InsertOneResult, error) {
	order.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.InsertOne(ctx, order)
	if err == nil {
		// Cache the new order
		orderJson, _ := json.Marshal(order)
		s.redisClient.Set(ctx, order.ID.Hex(), orderJson, 0)
	}
	return result, err
}

func (s *OrderService) GetOrderById(id primitive.ObjectID) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check Redis cache first
	val, err := s.redisClient.Get(ctx, id.Hex()).Result()
	if err == redis.Nil {
		// Cache miss, fetch from MongoDB
		var order models.Order
		err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
		if err != nil {
			return nil, err
		}

		// Cache the result in Redis
		orderJson, _ := json.Marshal(order)
		s.redisClient.Set(ctx, id.Hex(), orderJson, 0)
		return &order, nil
	} else if err != nil {
		return nil, err
	}

	// Cache hit, return cached order
	var order models.Order
	json.Unmarshal([]byte(val), &order)
	return &order, nil
}

func (s *OrderService) UpdateOrder(id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err == nil {
		// Invalidate cache
		s.redisClient.Del(ctx, id.Hex())
	}
	return result, err
}

func (s *OrderService) DeleteOrder(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err == nil {
		// Remove from cache
		s.redisClient.Del(ctx, id.Hex())
	}
	return result, err
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	// Find all documents in the collection
	cursor, err := s.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode each document into the Order struct
	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) GetOrderStatistics() ([]OrderStatistics, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "month", Value: bson.D{{Key: "$substr", Value: bson.A{"$createdAt", 0, 7}}}},
				{Key: "status", Value: "$status"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "_id.month", Value: 1},
			{Key: "_id.status", Value: 1},
		}}},
	}

	// Execute the aggregation
	cursor, err := s.collection.Aggregate(context.Background(), pipeline, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []OrderStatistics
	for cursor.Next(context.Background()) {
		var result OrderStatistics
		raw := cursor.Current

		// Extract and manually assign the values
		month := raw.Lookup("_id").Document().Lookup("month").StringValue()
		status := raw.Lookup("_id").Document().Lookup("status").StringValue()
		count := int(raw.Lookup("count").Int32())

		result.Month = month
		result.Status = status
		result.Count = count

		results = append(results, result)
	}

	return results, nil
}