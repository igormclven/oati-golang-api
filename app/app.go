package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	router = gin.Default()
)

type Student struct {
	Subject string `json:"subject"`
	Grade   int    `json:"grade"`
}

func loadConfig() {
	envFile, _ := godotenv.Read(".env")
	for key, value := range envFile {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}

}
func connectDB() (*mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func saveStudentGrade(c *gin.Context) {
	var student Student
	err := c.BindJSON(&student)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error while binding the request",
		})
		return
	}

	client, err := connectDB()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while connecting to the database",
		})
		return
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, context.Background())

	collection := client.Database("oati-api").Collection("students")
	_, err = collection.InsertOne(context.Background(), student)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while inserting the student",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "student saved successfully",
	})
}

func listAllGrades(c *gin.Context) {
	client, err := connectDB()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while connecting to the database",
		})
		return
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, context.Background())
	collection := client.Database("oati-api").Collection("students")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while finding the students",
		})
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, context.Background())

	var students []Student
	if err = cursor.All(context.Background(), &students); err != nil {
		c.JSON(500, gin.H{
			"error": "error while decoding the students",
		})
		return
	}
	c.JSON(200, students)
}

func Run() {
	loadConfig()
	router.POST("/saveGrade", saveStudentGrade)
	router.GET("/ListAllGrades", listAllGrades)

	err := router.Run(":8080")

	if err != nil {
		log.Fatal("Error while running the server: ", err)
	}
}
