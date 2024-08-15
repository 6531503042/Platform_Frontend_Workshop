package config

const (
	//Redis Config
	RedisAddr    = "localhost:6379"
	RedisPass	 = ""
	RedisDB		 = 0

	//MongoDB Config
	MongoURI     = "mongodb://localhost:27017"
	MongoUser	 = "admin"
	MongoPass	 = "admin"

	//Kafka Config
	KafkaBroker  = "localhost:9092"
	KafkaTopic   = "test_topic"

	//Server Config
	Port         = "8080"
)
