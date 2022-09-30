package models

import (
	"reps-hrms/pkg/config"
	"reps-hrms/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	ID          string      `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	Topics      []*Topic    `json:"topics"`
	Languages   []*Language `json:"languages"`

	Programmer *Programmer `json:"programmer"`
}

type Topic struct {
	Name string `json:"name"`
}

type Language struct {
	Name string `json:"name"`
}

type Programmer struct {
	Login string `json:"login"`
	Bio   string `json:"bio"`
	Email string `json:"email"`
}

const (
	REPS_COLL string = "repository"
)

var Mongo *config.MongoDB

func init() {
	var err error
	Mongo, err = config.InitMongoDB()
	if err != nil {
		utils.ErrLogger.Fatalf("Failed to initialize to MongoDB - %s\n", err.Error())
	}
}

func GetRepositories(c *fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := Mongo.DB.Collection(REPS_COLL).Find(c.Context(), query)
	if err != nil {
		utils.ErrLogger.Printf("Failed to query all repositories from MongoDB - %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	var reps []Repository = make([]Repository, 0)
	err = cursor.All(c.Context(), &reps)
	if err != nil {
		utils.ErrLogger.Printf("Failed to extract all repositories from MongoDB Query result- %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	utils.InfoLogger.Printf("Sent %d Repository records records ('GET':'/programmer')\n", len(reps))
	return c.Status(200).JSON(reps)
}

func CreateRepository(c *fiber.Ctx) error {
	collection := Mongo.DB.Collection(REPS_COLL)

	repo := &Repository{}
	err := c.BodyParser(repo)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request body into a Repository object - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	repo.ID = ""
	result, err := collection.InsertOne(c.Context(), repo)
	if err != nil {
		utils.ErrLogger.Printf("Failed to insert a Repository record to MondoDB - %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRepo := &Repository{}

	err = collection.FindOne(c.Context(), query).Decode(createdRepo)
	if err != nil {
		utils.ErrLogger.Printf("Failed to decode newly created Repository record to MondoDB - %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	utils.InfoLogger.Printf("Created a new Repository record ('POST':'/repository')\n")
	return c.Status(200).JSON(createdRepo)
}

func UploadRepositories(c *fiber.Ctx) error {
	collection := Mongo.DB.Collection(REPS_COLL)

	var repos []Repository
	err := c.BodyParser(&repos)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request body into a Repository objects - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	reposInterface := make([]interface{}, len(repos))
	for i, repo := range repos {
		repo.ID = ""
		reposInterface[i] = repo
	}

	result, err := collection.InsertMany(c.Context(), reposInterface)
	if err != nil {
		utils.ErrLogger.Printf("Failed to insert Repository records to MondoDB - %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	utils.InfoLogger.Printf("Created %d new Repository records ('POST':'/repository/upload')\n", len(result.InsertedIDs))
	return c.Status(200).JSON(result.InsertedIDs)
}

func GetRepositoryByID(c *fiber.Ctx) error {
	param := c.Params("id")
	repID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request param into a Repository ID - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	query := bson.D{{Key: "_id", Value: repID}}
	repo := &Repository{}

	err = Mongo.DB.Collection(REPS_COLL).FindOne(c.Context(), query).Decode(repo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrLogger.Printf("A Repository record with ID(%s) does not exist\n", param)
			return c.Status(404).JSON(utils.NotFoundMsg)
		}
		utils.ErrLogger.Printf("Failed to decode newly created Repository record to MondoDB - %s\n", err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	utils.InfoLogger.Printf("Sent a Repository record ('GET':'/repository/%s')\n", param)
	return c.Status(200).JSON(repo)
}

func DeleteRepository(c *fiber.Ctx) error {
	param := c.Params("id")
	repID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request param into a Repository ID - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	query := bson.D{{Key: "_id", Value: repID}}
	result, err := Mongo.DB.Collection(REPS_COLL).DeleteOne(c.Context(), query)
	if err != nil {
		utils.ErrLogger.Printf("Failed to delete a Repository record with ID(%s) - %s\n", param, err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}
	if result.DeletedCount < 1 {
		utils.ErrLogger.Printf("A Repository record with ID(%s) does not exist\n", param)
		return c.Status(404).JSON(utils.NotFoundMsg)
	}

	utils.InfoLogger.Printf("Deleted a Repository record ('DELETE':'/repository/%s')\n", param)
	return c.Status(200).JSON(utils.OKMsg)
}

func UpdateRepositoryByID(c *fiber.Ctx) error {
	param := c.Params("id")
	repID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request param into a Repository ID - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	repo := &Repository{}
	err = c.BodyParser(repo)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse request body into a Repository object - %s\n", err.Error())
		return c.Status(400).JSON(utils.BadRequestMsg)
	}

	query := bson.D{{Key: "_id", Value: repID}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: repo.Name},
				{Key: "description", Value: repo.Description},
				{Key: "license", Value: repo.License},
			},
		},
	}

	err = Mongo.DB.Collection(REPS_COLL).FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrLogger.Printf("A Repository record with ID(%s) does not exist\n", param)
			return c.Status(404).JSON(utils.NotFoundMsg)
		}
		utils.ErrLogger.Printf("Failed to update a Repository record with ID(%s) - %s\n", param, err.Error())
		return c.Status(500).JSON(utils.ServerErrorMsg)
	}

	repo.ID = param
	utils.InfoLogger.Printf("Updated a Repository record ('PUT':'/repository/%s')\n", param)
	return c.Status(200).JSON(repo)
}
