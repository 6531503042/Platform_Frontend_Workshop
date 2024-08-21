package services

import (
	"backend/models"
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	collection  *mongo.Collection
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
		collection:  collection,
		redisClient: redisClient,
	}
}

func (s *ProductService) CreateProduct(product models.Product) (*mongo.InsertOneResult, error) {
	// Validate required fields
	if product.Name == "" {
		return nil, errors.New("missing required field: name")
	}
	if product.Price == 0 {
		return nil, errors.New("missing required field: price")
	}

	// Insert product into MongoDB
	result, err := s.collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, errors.New("failed to insert product into MongoDB: " + err.Error())
	}

	// Cache the product in Redis
	productData, err := json.Marshal(product)
	if err != nil {
		return nil, errors.New("failed to marshal product data: " + err.Error())
	}

	cacheKey := "product:" + product.ID.Hex()
	err = s.redisClient.Set(context.Background(), cacheKey, productData, 0).Err()
	if err != nil {
		return nil, errors.New("failed to cache product in Redis: " + err.Error())
	}

	return result, nil
}

func (s *ProductService) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	// Check Redis cache
	cacheKey := "product:" + id.Hex()
	cachedProduct, err := s.redisClient.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		// Product not found in Redis, check MongoDB
		var product models.Product
		err = s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("product not found")
			}
			return nil, errors.New("failed to fetch product from MongoDB: " + err.Error())
		}

		// Cache result in Redis
		productData, err := json.Marshal(product)
		if err != nil {
			return nil, errors.New("failed to marshal product data: " + err.Error())
		}
		s.redisClient.Set(context.Background(), cacheKey, productData, 0)

		return &product, nil
	} else if err != nil {
		return nil, errors.New("failed to retrieve product from Redis: " + err.Error())
	}

	var product models.Product
	if err := json.Unmarshal([]byte(cachedProduct), &product); err != nil {
		return nil, errors.New("failed to unmarshal cached product data: " + err.Error())
	}

	return &product, nil
}

func (s *ProductService) ListProduct() ([]models.Product, error) {
	cursor, err := s.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, errors.New("failed to fetch products from MongoDB: " + err.Error())
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		return nil, errors.New("failed to decode products: " + err.Error())
	}
	return products, nil
}

func (s *ProductService) UpdateProduct(id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}
	result, err := s.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, errors.New("failed to update product in MongoDB: " + err.Error())
	}

	// Invalidate the cache
	cacheKey := "product:" + id.Hex()
	if err := s.redisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		return nil, errors.New("failed to invalidate cache: " + err.Error())
	}

	return result, nil
}

func (s *ProductService) DeleteProduct(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	cacheKey := "product:" + id.Hex()

	// Delete product from MongoDB
	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, errors.New("failed to delete product from MongoDB: " + err.Error())
	}

	// Invalidate the cache
	if err := s.redisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		return nil, errors.New("failed to invalidate cache: " + err.Error())
	}

	return result, nil
}

func (s *ProductService) GetProductCount() (int64, error) {
	count, err := s.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return 0, errors.New("failed to count products in MongoDB: " + err.Error())
	}
	return count, nil
}

func (s *ProductService) AggregateProducts(pipeline mongo.Pipeline) ([]ProductStatistics, error) {
	cursor, err := s.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, errors.New("failed to aggregate products in MongoDB: " + err.Error())
	}
	defer cursor.Close(context.Background())

	var statistics []ProductStatistics
	if err := cursor.All(context.Background(), &statistics); err != nil {
		return nil, errors.New("failed to decode aggregation results: " + err.Error())
	}

	return statistics, nil
}
