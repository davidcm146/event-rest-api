package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidcm146/event-rest-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

// registerUser godoc
// @Summary Register a new user
// @Description Create a new user account with email, password and name
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   register body RegisterUserRequest true "User registration payload"
// @Success 201 {object} map[string]interface{} "Successfully registered"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /register [post]
func (app *application) registerUser(c *gin.Context) {
	var register RegisterUserRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	register.Password = string(hashedPassword)
	user := &database.User{
		Email:    register.Email,
		Name:     register.Name,
		Password: register.Password,
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// loginUser godoc
// @Summary Login user
// @Description Authenticate user with email and password, returns a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   login body LoginUserRequest true "User login payload"
// @Success 200 {object} LoginUserResponse "Successfully authenticated"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /login [post]
func (app *application) loginUser(c *gin.Context) {
	var auth LoginUserRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	fmt.Println("User:", err)
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"expire": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, LoginUserResponse{Token: tokenString})
}
