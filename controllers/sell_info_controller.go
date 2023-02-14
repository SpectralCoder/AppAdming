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

var sellInfoCollection *mongo.Collection = configs.GetCollection(configs.DB, "sells")
var sellValidate = validator.New()

func CreateSell() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var sell models.SellInfo
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&sell); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := sellValidate.Struct(&sell); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newSell := models.SellInfo{
			Id:              primitive.NewObjectID(),
			Products:        sell.Products,
			Customer_id:     sell.Customer_id,
			Organization_id: sell.Organization_id,
			Amount:          sell.Amount,
		}

		result, err := sellInfoCollection.InsertOne(ctx, newSell)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CommonResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetASell() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		sellId := c.Param("sellId")
		var sell models.SellInfo
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(sellId)

		err := sellInfoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&sell)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": sell}})
	}
}

func EditASell() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		sellId := c.Param("sellId")
		var sell models.SellInfo
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(sellId)

		//validate the request body
		if err := c.BindJSON(&sell); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := sellValidate.Struct(&sell); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"Products":        sell.Products,
			"Customer_id":     sell.Customer_id,
			"Organization_id": sell.Organization_id,
			"Amount":          sell.Amount,
		}
		result, err := sellInfoCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated sell details
		var updatedSell models.SellInfo
		if result.MatchedCount == 1 {
			err := sellInfoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedSell)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedSell}})
	}
}

func DeleteASell() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		sellId := c.Param("sellId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(sellId)

		result, err := sellInfoCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.CommonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "sell with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "sell successfully deleted!"}},
		)
	}
}

func GetAllSells() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var sells []models.SellInfo
		defer cancel()

		results, err := sellInfoCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleSell models.SellInfo
			if err = results.Decode(&singleSell); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			sells = append(sells, singleSell)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": sells}},
		)
	}
}
