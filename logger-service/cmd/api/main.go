package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.cm/obynonwane/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:2701"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	//Context is used in managing the lifecycle of an operation - fundamental to Go
	//create a context in other to disconnect after 15 seconds of any database operation
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//this line schedules a call for the context to exit if the main function exits - sorrounding function
	//it is essential to cancel when the sorrounding function exit so as to cleanup resources
	defer cancel() //always call to release resource

	//another defer statement
	//that defines an anonymous function
	//that disconnects from the db when the sorrounding function
	// in this case main exits
	defer func() {
		//Attenpts disconnecting
		//checks if there is an error during disconnection and then panics
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//start web server
	go app.serve()
}

// create a web server function
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	//create connection object
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//make the connection
	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panic(err)
	}

	return c, nil
}
