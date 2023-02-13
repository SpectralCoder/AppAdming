package controllers

import (
	"appadming/configs"
	"appadming/models"
	"appadming/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var customerCollection *mongo.Collection = configs.GetCollection(configs.DB, "customers")
var validate = validator.New()

func CreateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var customer models.Customer
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&customer); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newCustomer := models.Customer{
			Id:       primitive.NewObjectID(),
			Name:     customer.Name,
			Father:   customer.Father,
			Home:     customer.Home,
			Village:  customer.Village,
			Thana:    customer.Thana,
			District: customer.District,
			Paid:     customer.Paid,
			Due:      customer.Due,
			Phone:    customer.Phone,
			Email:    customer.Email,
		}

		result, err := customerCollection.InsertOne(ctx, newCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CommonResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetACustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		customerId := c.Param("customerId")
		var customer models.Customer
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(customerId)

		err := customerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": customer}})
	}
}

func EditACustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		customerId := c.Param("customerId")
		var customer models.Customer
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(customerId)

		//validate the request body
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&customer); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"name":     customer.Name,
			"father":   customer.Father,
			"home":     customer.Home,
			"village":  customer.Village,
			"thana":    customer.Thana,
			"district": customer.District,
			"paid":     customer.Paid,
			"due":      customer.Due,
			"phone":    customer.Phone,
			"email":    customer.Email,
		}
		result, err := customerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated customer details
		var updatedCustomer models.Customer
		if result.MatchedCount == 1 {
			err := customerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedCustomer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedCustomer}})
	}
}

func DeleteACustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		customerId := c.Param("customerId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(customerId)

		result, err := customerCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.CommonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "customer with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "customer successfully deleted!"}},
		)
	}
}

func GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var customers []models.Customer
		defer cancel()

		results, err := customerCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleCustomer models.Customer
			if err = results.Decode(&singleCustomer); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			customers = append(customers, singleCustomer)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": customers}},
		)
	}
}
