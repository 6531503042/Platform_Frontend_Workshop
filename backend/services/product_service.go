package services

import (
	"backend/models"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	collection *mongo.Collection
	redisClient *redis.Client
}

type ProductStatistics struct {
	Month  string `json:"month"`
	Count  int    `json:"count"`
	Status string `json:"status,omitempty"`
}

// NewProductService creates a new instance of ProductService
func NewProductService(collection *mongo.Collection, redisClient *redis.Client) *ProductService {
	return &ProductService{
		collection: collection,
		redisClient: redisClient,
	}
}

func (s *ProductService) CreateProduct(product models.Product) (*mongo.InsertOneResult, error) {
	result, err := s.collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}

	// Cache result in Redis
	productData, _ := json.Marshal(product)
	cacheKey := "product:" + product.ID.Hex()
	s.redisClient.Set(context.Background(), cacheKey, productData, 0)

	return result, nil
}

func (s *ProductService) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	// Check Redis cache
	cacheKey := "product:" + id.Hex()
	cachedProduct, err := s.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var product models.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err != nil {
			return nil, err
		}
		return &product, nil
	}

	// Fetch from MongoDB
	var product models.Product
	err = s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}

	// Cache result in Redis
	productData, _ := json.Marshal(product)
	s.redisClient.Set(context.Background(), cacheKey, productData, 0)

	return &product, nil
}

func (s *ProductService) ListProduct() ([]models.Product, error) {
	cursor, err := s.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) UpdateProduct(id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}
	result, err := s.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Invalidate the cache
	cacheKey := "product:" + id.Hex()
	s.redisClient.Del(context.Background(), cacheKey)

	return result, nil
}

func (s *ProductService) DeleteProduct(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	cacheKey := "product:" + id.Hex()

	// Delete product from MongoDB
	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	// Invalidate the cache
	s.redisClient.Del(context.Background(), cacheKey)

	return result, nil
}

func (s *ProductService) GetProductCount() (int64, error) {
	count, err := s.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ProductService) AggregateProducts(pipeline mongo.Pipeline) ([]ProductStatistics, error) {
	cursor, err := s.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var statistics []ProductStatistics
	if err := cursor.All(context.Background(), &statistics); err != nil {
		return nil, err
	}

	return statistics, nil
}

