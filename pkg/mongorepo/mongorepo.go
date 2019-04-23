package mongorepo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Repository struct {
	*mongo.Client
	DatabaseName string
}

// "mongodb://localhost:27017"
func NewRepository(host string, dbName string) *Repository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(host))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return &Repository{client, dbName}
}

func (r *Repository) Create(ctx context.Context, collectionName string, doc interface{}) (*mongo.InsertOneResult, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

func (r *Repository) Read(ctx context.Context, collectionName string, filter [][]string) (*mongo.Cursor, error) {
	database := r.Database(r.DatabaseName)
	collection := database.Collection(collectionName)
	bsonFilter := r.buildFilter(filter)
	res, err := collection.Find(ctx, bsonFilter)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

func (r *Repository) buildFilter(filter [][]string) bson.M {
	bsonFilter := bson.M{}
	for _, f := range filter {
		bsonFilter[f[0]] = bson.M{"$elemMatch": bson.M{"$eq": f[1]}}
	}
	return bsonFilter
}
