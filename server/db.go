package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConnection Struct for db configuration, connection, collection and context
type DBConnection struct {
	Host       string
	Port       int
	DB         *mongo.Client
	Collection *mongo.Collection
	MongoCtx   context.Context
}

// Global variable for DBConnection
var dbConn DBConnection

func (dbConn *DBConnection) connectionString() string {
	return fmt.Sprintf("mongodb://%s:%d", dbConn.Host, dbConn.Port)
}

// Connect tries to connect to MongoDB.
func (dbConn *DBConnection) Connect() (err error) {
	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")

	// non-nil empty context
	dbConn.MongoCtx = context.Background()

	// Connect takes in a context and options, the connection URI is the only option we pass for now
	dbConn.DB, err = mongo.Connect(dbConn.MongoCtx, options.Client().ApplyURI(dbConn.connectionString()))
	// Handle potential errors
	if err != nil {
		log.Fatal(err)
	}

	// Check whether the connection was succesful by pinging the MongoDB server
	err = dbConn.DB.Ping(dbConn.MongoCtx, nil)
	return
}

// SetCollection returns formatted Connection String.
func (dbConn *DBConnection) SetCollection(database, collection string) {
	// Bind our collection to our global variable for use in other methods
	dbConn.Collection = dbConn.DB.Database(database).Collection(collection)
}

// Disconnect do transparent disconnect.
func (dbConn *DBConnection) Disconnect() {
	dbConn.DB.Disconnect(dbConn.MongoCtx)
}
