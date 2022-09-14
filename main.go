package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipes-api/handlers"
	"recipes-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	recipes []models.Recipe
	ctx     context.Context
	err     error
	client  *mongo.Client
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx = context.Background()

	// Config mongodb
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	// Config redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)

	router.Run()
}
