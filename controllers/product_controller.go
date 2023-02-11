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

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")
var productValidate = validator.New()

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := productValidate.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newProduct := models.Product{
			Id:           primitive.NewObjectID(),
			Model:        product.Model,
			Price:        product.Price,
			Cost:         product.Cost,
			Description:  product.Description,
			Category:     product.Category,
			ImageURL:     product.ImageURL,
			Stock:        product.Stock,
			Brand:        product.Brand,
			Organization: product.Organization,
		}

		result, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.CommonResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": product}})
	}
}

func EditAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(productId)

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := productValidate.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.CommonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"category":     product.Category,
			"brand":        product.Brand,
			"cost":         product.Cost,
			"description":  product.Description,
			"image_url":    product.ImageURL,
			"model":        product.Model,
			"price":        product.Price,
			"stock":        product.Stock,
			"organization": product.Organization,
		}
		result, err := productCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated product details
		var updatedProduct models.Product
		if result.MatchedCount == 1 {
			err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedProduct}})
	}
}

func DeleteAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		productId := c.Param("productId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		result, err := productCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.CommonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "product with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "product successfully deleted!"}},
		)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var products []models.Product
		defer cancel()

		results, err := productCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, responses.CommonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			products = append(products, singleProduct)
		}

		c.JSON(http.StatusOK,
			responses.CommonResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": products}},
		)
	}
}
