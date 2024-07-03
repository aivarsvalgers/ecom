package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aivarsvalgers/ecom/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection, *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		UserCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.Handlerfunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")

			_ = c.AbortWithError(http.StatusBadRequest, error.New("Product ID is empty!"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User ID is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.UserCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully added to the cart!")
	}
}

func (app *Application) RemoveItem() gin.Handlerfunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")

			_ = c.AbortWithError(http.StatusBadRequest, error.New("Product ID is empty!"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User ID is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.UserCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully removed item from cart!")
	}
}

func GetItemFromCart() gin.Handlerfunc {
	return func(c *gin.Context) {
		c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid id!"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcard models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive{Key:"_id", Value: usert_id}}.Decode(&filledcard))

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Not found!")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive{Key "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive{Key:"$sum", Value: "$usercart.price"}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline(filter_match, unwind, grouping))

		if err != nil {
			log.Println(err)
		}
		
		var listing {}bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _. json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}

		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.Handlerfunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panicln("User ID is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty!"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.UserCollection, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON("Successfully placed the order!")
	}
}

func (app *Application) InstantBuy() gin.Handlerfunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")

			_ = c.AbortWithError(http.StatusBadRequest, error.New("Product ID is empty!"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User ID is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.UserCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully placed the order!")
	}
}
