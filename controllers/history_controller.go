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

var historyCollection *mongo.Collection = configs.GetCollection(configs.DB, "historys")
var historyValidate = validator.New()

func CreateHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var history models.History
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&history); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := historyValidate.Struct(&history); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newHistory := models.History{
			Id:          primitive.NewObjectID(),
			Due:         history.Due,
			Paid:        history.Paid,
			Date:        history.Date,
			Customer_id: history.Customer_id,
			Seller_id:   history.Seller_id,
		}

		result, err := historyCollection.InsertOne(ctx, newHistory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CommonResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		historyId := c.Param("historyId")
		var history models.History
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(historyId)

		err := historyCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&history)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": history}})
	}
}

func EditAHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		historyId := c.Param("historyId")
		var history models.History
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(historyId)

		//validate the request body
		if err := c.BindJSON(&history); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := historyValidate.Struct(&history); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"Due":         history.Due,
			"Paid":        history.Paid,
			"Date":        history.Date,
			"Customer_id": history.Customer_id,
			"seller_id":   history.Seller_id,
		}
		result, err := historyCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated history details
		var updatedHistory models.History
		if result.MatchedCount == 1 {
			err := historyCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedHistory)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedHistory}})
	}
}

func DeleteAHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		historyId := c.Param("historyId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(historyId)

		result, err := historyCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.CommonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "history with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "history successfully deleted!"}},
		)
	}
}

func GetAllHistorys() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var historys []models.History
		defer cancel()

		results, err := historyCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleHistory models.History
			if err = results.Decode(&singleHistory); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			historys = append(historys, singleHistory)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": historys}},
		)
	}
}
