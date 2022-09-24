package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type Todo struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
    Title string  	`json:"title"`
    Des  string 		`json:"des"`
    Completed bool 	`json:"completed"`
}
func main() {
	router := gin.Default()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel();
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { 
		log.Fatal(err)
	}
	collection := client.Database("golang").Collection("todos")

	//Get all todos
	router.GET("/todos", func(c *gin.Context){
		cur, err := collection.Find(context.Background(), bson.D{})
		if err != nil { log.Fatal(err) }
		defer cur.Close(context.Background())

		var results []Todo
		if err = cur.All(context.Background(), &results); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"todos": results,
		})
	})

	//Create new todo
	router.POST("/todos", func(c *gin.Context){
		var newTodo Todo
		if err := c.BindJSON(&newTodo); err != nil {
			return
		}
		res, err := collection.InsertOne(context.Background(), newTodo)
		if err != nil { 
			log.Fatal(err)
		 }
		id := res.InsertedID
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"_id": id,
		})
	})

	//Get todo by id
	router.GET("/todos/:id", func(c *gin.Context){
		var id = c.Param("id")
		todoIdObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"message": "Id invalid",
			})
			return
		}
		var result Todo
		err = collection.FindOne(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)

		if result.Title == "" && result.Des == "" {
				c.JSON(http.StatusNotFound, gin.H{
				"status": "failure",
				"message": "Todo not found!",
			})
			return
		}else{
				c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"item": result,
			})
			return
		}
	})
	router.DELETE("/todos/:id", func(c *gin.Context){
		var id = c.Param("id")
		todoIdObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"message": "Id invalid",
			})
			return
		}
		var result Todo
		err = collection.FindOneAndDelete(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)
		if result.Title == "" && result.Des == "" {
				c.JSON(http.StatusNotFound, gin.H{
				"status": "failure",
				"message": "Todo not found!",
			})
			return
		}else{
				c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"message": "Deleted todo successfully!",
			})
			return
		}
	})
	router.Run("localhost:3000")
}