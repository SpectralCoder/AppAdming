package controllers

import (
	"appadming/configs"
	"appadming/models"
	"appadming/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = configs.GetCollection(configs.DB, "organizations")
var organizationValidate = validator.New()

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var organization models.Orgnization
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&organization); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		uidInterface, _ := c.Get("uid")
		uid, _ := uidInterface.(string)

		uObjectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}

		//use the validator library to validate required fields
		if validationErr := organizationValidate.Struct(&organization); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		newOrganization := models.Orgnization{
			Id:         primitive.NewObjectID(),
			Name:       organization.Name,
			Address:    organization.Address,
			Phone:      organization.Phone,
			Email:      organization.Email,
			Created_at: time.Time{},
		}
		newOrganization.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := organizationCollection.InsertOne(ctx, newOrganization)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		_, err = processCreateUserOrganization(ctx, newOrganization.Id, uObjectID, "active", "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CommonResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		organizationId := c.Param("organizationId")
		var organization models.Orgnization
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(organizationId)

		err := organizationCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&organization)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": organization}})
	}
}

func EditAOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		organizationId := c.Param("organizationId")
		var organization models.Orgnization
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(organizationId)

		//validate the request body
		if err := c.BindJSON(&organization); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := organizationValidate.Struct(&organization); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"Id":      organization.Id,
			"Name":    organization.Name,
			"Address": organization.Address,
			"Phone":   organization.Phone,
			"Email":   organization.Email,
		}
		result, err := organizationCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated organization details
		var updatedOrganization models.Orgnization
		if result.MatchedCount == 1 {
			err := organizationCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedOrganization)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedOrganization}})
	}
}

func DeleteAOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		organizationId := c.Param("organizationId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(organizationId)

		result, err := organizationCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.CommonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "organization with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "organization successfully deleted!"}},
		)
	}
}

func GetAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var organizations []models.Orgnization
		defer cancel()

		results, err := organizationCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleOrganization models.Orgnization
			if err = results.Decode(&singleOrganization); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			organizations = append(organizations, singleOrganization)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": organizations}},
		)
	}
}

func GetSearchOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		search := c.Query("name")
		var organizations []models.Orgnization
		defer cancel()

		results, err := organizationCollection.Find(ctx, bson.M{"name": bson.M{"$regex": search, "$options": "i"}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleOrganization models.Orgnization
			if err = results.Decode(&singleOrganization); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			organizations = append(organizations, singleOrganization)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": organizations}},
		)
	}
}
