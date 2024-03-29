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
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"title": title}
	update := bson.M{
			"$set": bson.D{
				{"body", body},
				{"updated_at", timestamp}},
			"$setOnInsert": bson.D{
				{"created_at", timestamp}}}
	db.collection.UpdateOne(ctx, filter, update, opts)
}
