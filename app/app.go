package app

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Student struct {
	Subject string `json:"subject"`
	Grade   int    `json:"grade"`
}

type Subject struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var (
	router = gin.Default()
)

// Loads the config from the .env file
func loadConfig() {
	envFile, _ := godotenv.Read(".env")
	for key, value := range envFile {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}

}

// Connects to the database
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
	subjectCollection := client.Database("oati-api").Collection("subjects")

	var subject Subject
	err = subjectCollection.FindOne(context.Background(), bson.M{"name": student.Subject}).Decode(&subject)
	print(subject.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "subject code not found",
		})
		return
	}

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

func saveSubject(c *gin.Context) {
	var subject Subject
	err := c.BindJSON(&subject)
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

	var existingSubject Subject

	collection := client.Database("oati-api").Collection("subjects")

	err = collection.FindOne(context.Background(), bson.M{"code": subject.Code}).Decode(&existingSubject)
	if err != nil {
		_, err = collection.InsertOne(context.Background(), subject)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "error while inserting the subject",
			})
			return
		}
	} else {
		c.JSON(400, gin.H{
			"error": "subject already exists",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "subject saved successfully",
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
	// Load the config from the .env file
	loadConfig()

	// Enable CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	// Routes
	router.POST("/saveGrade", saveStudentGrade)
	router.POST("/saveSubject", saveSubject)
	router.GET("/listAllGrades", listAllGrades)

	err := router.Run(":8080")

	if err != nil {
		log.Fatal("Error while running the server: ", err)
	}
}
