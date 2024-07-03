package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aivarsvalgers/ecom/database"
	"github.com/aivarsvalgers/ecom/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(databsase.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword (password string) string {
	bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword (userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login password is incorrect!"
		valid = false
	}

	return valid, msg
}

func Signup() gin.Handlerfunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatonErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists!"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerErro, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This phone no. is already in use!"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.Fist_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The user didn't get created!"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.Handlerfunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email:" user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login or password incorrect!"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		
		defer cancel()

		if !PasswordIsValid {
			c.JSON{http.StatusInternalServerError, gin.H{"error": msg}}
			fmt.Println(msg)
			return
		}
		token, regreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, *founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}
}

func ProductViewerAdmin() gin.Handlerfunc {

}

func SearchProduct() gin.Handlerfunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, please try after some time!")
			return
		}

		err = cursor.All(ctx, &productlist)

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}

		defer cursor.Close()

		if err := cursor.err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid!")
			return
		}
		defer cancel()
		
		c.IndentedJSON(200, productlist)
	}
}

func SearchProductByQuery() gin.Handlerfunc {
	return func(c *gin.Context) {
		var SearchProducts []models.Product

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("Query is empty!")
			c.Header("Content-Type", "application/json")
			c.JSON(HTTP.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})

		if err != nil {
			c.IndentedJSON(404, "Something went wrong while fetching the data!")
			return
		}

		err = searchquerydb.All(ctx, &searchproducts)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid!")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invlaid request!")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}
