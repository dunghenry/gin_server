package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title"`
	Des       string             `json:"des"`
	Completed bool               `json:"completed"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// production
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("golang").Collection("todos")
	fmt.Println("Connect MongoDB")
	//Get all todos
	router.GET("/todos", func(c *gin.Context) {
		cur, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(context.Background())

		var results []Todo
		if err = cur.All(context.Background(), &results); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"todos":  results,
		})
	})

	//Create new todo
	router.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		newTodo.Id = primitive.NewObjectID()
		if err := c.BindJSON(&newTodo); err != nil {
			return
		}
		res, err := collection.InsertOne(context.Background(), newTodo)
		if err != nil {
			log.Fatal(err)
		}
		id := res.InsertedID
		fmt.Println(res)
		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"_id":    id,
		})
	})

	//Get todo by id
	router.GET("/todos/:id", func(c *gin.Context) {
		var id = c.Param("id")
		todoIdObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "failure",
				"message": "Id invalid",
			})
			return
		}
		var result Todo
		err = collection.FindOne(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)
		if result.Title == "" && result.Des == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failure",
				"message": "Todo not found!",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"item":   result,
			})
			return
		}
	})
	router.DELETE("/todos/:id", func(c *gin.Context) {
		var id = c.Param("id")
		todoIdObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "failure",
				"message": "Id invalid",
			})
			return
		}
		var result Todo
		err = collection.FindOneAndDelete(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)
		if result.Title == "" && result.Des == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failure",
				"message": "Todo not found!",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Deleted todo successfully!",
			})
			return
		}
	})

	// Update todo by id
	router.PUT("todos/:id", func(c *gin.Context) {
		var id = c.Param("id")
		todoIdObject, err := primitive.ObjectIDFromHex(id)
		var updateTodo Todo
		if err := c.BindJSON(&updateTodo); err != nil {
			return
		}
		fmt.Println(updateTodo.Title)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failure",
				"message": "Id invalid",
			})
			return
		}
		var result Todo
		err = collection.FindOne(context.TODO(), bson.D{{"_id", todoIdObject}}).Decode(&result)
		if result.Title == "" && result.Des == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failure",
				"message": "Todo not found!",
			})
			return
		} else {
			opts := options.Update().SetUpsert(true)
			result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", todoIdObject}}, bson.D{{"$set", bson.D{{"title", updateTodo.Title}, {"des", updateTodo.Des}, {"completed", updateTodo.Completed}}}}, opts)
			if err != nil {
				log.Fatal(err)
			}
			if result.MatchedCount != 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
					"item":   "Updated todo successfully!",
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failure",
					"message": "Todo update failure!",
				})
				return
			}
		}
	})
	router.Run(":3000")
}
