package handlers

import (
	"fmt"
	"sangeet-server/db"
	"sangeet-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func AddAuthGroup(app *fiber.App) {
	authGroup := app.Group("/auth")

	authGroup.Get("/", getAllUsers)
	authGroup.Get("/me", getMe)
	authGroup.Get("/:id", getUser)
	authGroup.Post("/register", registerUser)
	authGroup.Post("/login", loginUser)
	authGroup.Get("/refresh", refreshLogin)
}

/*
	TODO:

login
set cookies
*/
func registerUser(c *fiber.Ctx) error {
	var reqUser models.User
	err := c.BodyParser(&reqUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Request", "message": err.Error()})
	}

	if reqUser.FirstName == "" || reqUser.LastName == "" || reqUser.Password == "" ||
		reqUser.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Please fill out all the fields"})
	}

	coll := db.GetDBCollection("users")

	// verify unique email
	userFound := coll.FindOne(c.Context(), bson.M{"user_email": reqUser.Email})
	if userFound.Err() == nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "User exists", "message": userFound.Err().Error()})
	}

	// pass length check
	if len(reqUser.Password) < 8 || len(reqUser.Password) > 20 {
		return c.Status(400).JSON(fiber.Map{"error": "Password should be between 8-20 characters."})
	}

	// hash pass
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), 8)
	if err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "Please choose another password.", "message": err.Error()})
	}

	// create user
	var createUser models.User
	createUser.Email = reqUser.Email
	createUser.Password = string(hashedPass)
	createUser.FirstName = reqUser.FirstName
	createUser.LastName = reqUser.LastName

	// first user ? admin : user
	docCount, err := coll.CountDocuments(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Something went wrong"})
	}
	if docCount == 0 {
		createUser.Role = "admin"
	} else {
		createUser.Role = "user"
	}

	createdUser, err := coll.InsertOne(c.Context(), createUser)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Something went wrong", "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"res": createdUser,
	})
}

/*
	TODO:

jwt access and refresh token
set cookies
*/
func loginUser(c *fiber.Ctx) error {
	var reqUser models.User
	err := c.BodyParser(&reqUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Request", "message": err.Error()})
	}

	if reqUser.Password == "" || reqUser.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Please fill out all the fields"})
	}
	if len(reqUser.Password) < 8 || len(reqUser.Password) > 20 {
		return c.Status(400).JSON(fiber.Map{"error": "Password should be between 8-20 characters."})
	}

	coll := db.GetDBCollection("users")

	// verify user exists
	var userFound models.User
	err = coll.FindOne(c.Context(), bson.M{"user_email": reqUser.Email}).Decode(&userFound)
	if err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "We could not find you", "message": err.Error()})
	}

	// verify pass
	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(reqUser.Password))
	if err != nil {
		return c.Status(400).
			JSON(fiber.Map{"error": "Incorrect password", "message": err.Error()})
	}

	token, err := GetToken(c, models.Token{UserRole: userFound.Role, UserEmail: userFound.Email})
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{"error": "Something went wrong", "message": err.Error()})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "access-token"
	cookie.Value = token
	cookie.Expires.Add(72 * time.Hour).Unix()
	cookie.HTTPOnly = true

	c.Cookie(cookie)
	fmt.Println("set cookie", cookie)
	fmt.Println("token:", token)

	return c.SendStatus(200)
}

func refreshLogin(c *fiber.Ctx) error {
	return nil
}

/*
	TODO:

read cookies
get user
*/
func getMe(c *fiber.Ctx) error {
	return nil
}

/*
NOTE:
only for admin
*/
func getAllUsers(c *fiber.Ctx) error {
	return nil
}

/*
NOTE:
only for admin
*/
func getUser(c *fiber.Ctx) error {
	return nil
}
