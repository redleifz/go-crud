// controllers/user_controller.go

package controllers

import (
	"database/sql"
	"fmt"
	"go-crud/configs"
	"go-crud/models"
	"go-crud/responses"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *gin.Context) {
	// Get the DB instance from your configs package
	db := configs.ConnectDB()

	// Query all users from the database
	rows, err := db.Query("SELECT user_id, user_login, user_pwd FROM user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserLogin, &user.UserPwd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a custom JSON structure
	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success",
		Data: map[string]interface{}{"results": users, "total": len(users)}})
}

func QueryUsers() (*sql.Rows, error) {
	db := configs.ConnectDB()
	rows, err := db.Query("SELECT user_id, user_login, user_pwd FROM user")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetAllUsersExcelFile(c *gin.Context) {
	// Query all users from the database
	rows, err := QueryUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Print the query result for debugging
	for rows.Next() {
		var userID int
		var userLogin string
		var userPwd string
		if err := rows.Scan(&userID, &userLogin, &userPwd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("%d,%s,%s\n", userID, userLogin, userPwd)
	}
}

// PostUser
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func UserLogin(c *gin.Context) {

	//match value from json body request
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := configs.ConnectDB()
	// var foundUser models.User
	var storedUsername, storedHashedPassword string

	// fmt.Println(user.UserLogin)

	// Query the database to get username and password
	err := db.QueryRow("SELECT user_login, user_pwd FROM user WHERE user_login = ?", user.UserLogin).Scan(&storedUsername, &storedHashedPassword)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.UserPwd))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	//jwt sign
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = storedUsername
	claims["exp"] = time.Now().AddDate(0, 0, 30).Unix() // Set expiration time to 30 days

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign the token"})
		return
	}

	// Authentication successful, return the token
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": tokenString})

	return
}

func CreateUser(c *gin.Context) {
	db := configs.ConnectDB()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user with the same login already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE user_login = ?", user.UserLogin).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		// A user with the same login already exists, send a response with a 409 Conflict status code
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "error": "User with the same login already exists"})
		return
	}

	hash, _ := HashPassword(user.UserPwd) // ignore error for the sake of simplicity

	fmt.Println("Password:", user.UserPwd)
	fmt.Println("Hash:    ", hash)

	// Insert user to the database
	_, err = db.Exec("INSERT INTO user (user_login, user_pwd, user_IDcard) VALUES (?, ?, ?)",
		user.UserLogin, hash, user.UserIdCard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success"})
}
