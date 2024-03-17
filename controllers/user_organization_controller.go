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

var userOrgCollection *mongo.Collection = configs.GetCollection(configs.DB, "user_organization_mapper")
var userOrgValidate = validator.New()

func processCreateUserOrganization(ctx context.Context, orgObjectID primitive.ObjectID, uObjectID primitive.ObjectID, status string, role string) (*mongo.InsertOneResult, error) {
	newUserOrg := models.UserOrganization{
		Id:           primitive.NewObjectID(),
		Organization: orgObjectID,
		User:         uObjectID,
		Role:         role,
		Status:       status,
	}

	result, err := userOrgCollection.InsertOne(ctx, newUserOrg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateUserOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := c.Param("org_id")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		uidInterface, _ := c.Get("uid")
		uid, _ := uidInterface.(string)

		var userOrg models.UserOrganizationRequest
		if err := c.BindJSON(&userOrg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgObjectID, err := primitive.ObjectIDFromHex(orgID)
		if err != nil {
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}

		uObjectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}
		//use the validator library to validate required fields
		if validationErr := userOrgValidate.Struct(&userOrg); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		result, err := processCreateUserOrganization(ctx, orgObjectID, uObjectID, "pending", userOrg.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"result": result})
	}
}

func GetAllUserOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		uidInterface, _ := c.Get("uid")
		uid, _ := uidInterface.(string)
		uObjectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			// Handle error (the hex string might not be a valid ObjectId)
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}
		var userOrgs []models.UserOrganization
		cursor, err := userOrgCollection.Find(ctx, bson.M{"user": uObjectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var userOrg models.UserOrganization
			cursor.Decode(&userOrg)
			userOrgs = append(userOrgs, userOrg)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": userOrgs}})
	}
}

func ApproveRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Param("id")
		id, err := primitive.ObjectIDFromHex(requestId)
		uidInterface, _ := c.Get("uid")
		uid, _ := uidInterface.(string)
		uObjectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			// Handle error (the hex string might not be a valid ObjectId)
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Fetch the organization ID associated with the request
		var requestDoc struct {
			Organization primitive.ObjectID `bson:"organization"`
		}
		err = userOrgCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&requestDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find request details"})
			return
		}

		orgObjectID := requestDoc.Organization

		filter := bson.M{
			"user":         uObjectID,
			"organization": orgObjectID,
			"role":         "admin",
			"status":       "active",
		}

		count, err := userOrgCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user role"})
			return
		}

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not an admin for the requested organization"})
			return
		}

		// If the user is an admin, update the request accordingly
		update := bson.M{
			"$set": bson.M{
				"status": "approved", // Set the status to "approved" or to your desired status
			},
		}

		_, err = userOrgCollection.UpdateByID(ctx, id, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Request approved successfully"})
	}

}

func GetUserOfOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		organizationId := c.Param("org_id")
		var userOrgs []models.UserOrganization
		objId, _ := primitive.ObjectIDFromHex(organizationId)
		cursor, err := userOrgCollection.Find(ctx, bson.M{"organization": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var userOrg models.UserOrganization
			cursor.Decode(&userOrg)
			userOrgs = append(userOrgs, userOrg)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": userOrgs}})
	}
}

func GetMyOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		uidInterface, _ := c.Get("uid")
		uid, _ := uidInterface.(string)
		uObjectID, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			// Handle error (the hex string might not be a valid ObjectId)
			fmt.Println("Error converting hex to ObjectId:", err)
			return
		}

		organizationId := c.Param("org_id")
		var userOrg models.UserOrganization
		objId, _ := primitive.ObjectIDFromHex(organizationId)
		cursor, err := userOrgCollection.Find(ctx, bson.M{"organization": objId, "user": uObjectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			cursor.Decode(&userOrg)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": userOrg}})
	}
}
