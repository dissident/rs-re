package db

import(
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dissident/rs-re/support"
)

type DB struct {
	client *mongo.Client
	collection *mongo.Collection
}

func InitDb(url string, dbName string, collectionName string) DB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	support.FailOnError(err, "ERROR with mongo connection")
	collection := client.Database(dbName).Collection(collectionName)
	return DB{ client, collection }
}

func (db *DB) Insert(title string, body string) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	timestamp := time.Now().Unix()
	newInsert := bson.D{
		{"title", title},
		{"body", body},
		{"created_at", timestamp}}
	db.collection.InsertOne(ctx, newInsert)
}
