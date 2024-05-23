package dbclient

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection attributes
const  (
	host = "localhost" 
	port = ":27017" 
	database = "services"
)

var MONGO_CLIENT *mongo.Client 
var DB_OBJ *mongo.Database

func ConnectMongo() {
	opt := options.Client().ApplyURI("mongodb://" + host + port)
	MONGO_CLIENT, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		log.Fatal("failed to establish a connection with the DB, err: ", err)
	}

	err = MONGO_CLIENT.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("DB irreachable, err: ", err)
	}
	DB_OBJ = MONGO_CLIENT.Database(database)
	
}

func DisconnectMongo() {
	err := MONGO_CLIENT.Disconnect(context.Background())
	if err != nil {
		log.Fatal("Error disconnecting from DB, err: ",  err)
	}
}