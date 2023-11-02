// controllers/user_controller.go

package controllers

import (
	"database/sql"
	"fmt"
	"go-crud/configs"
	"go-crud/models"
	"go-crud/responses"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func GetAllUsers(c *gin.Context) {
	// Get the DB instance from your configs package
	db := configs.ConnectDB()

	// Query all users from the database
	rows, err := db.Query("SELECT user_id, user_login, user_citizen_id, user_password FROM user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserLogin, &user.UserCitizenId, &user.UserPassword); err != nil {
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
	rows, err := db.Query("SELECT user_id, user_login, user_password FROM user")
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
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := configs.ConnectDB()
	var storedUsername, storedHashedPassword, storeRole string

	err := db.QueryRow("SELECT user_login, user_password ,user_role FROM user WHERE user_login = ?", user.UserLogin).Scan(&storedUsername, &storedHashedPassword, &storeRole)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.UserPassword))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// JWT signing process using the loaded secret key from environment variable
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	fmt.Println(`secret`, secretKey)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = storedUsername
	if storeRole == "3" {
		claims["role"] = "admin"
	} else {
		claims["role"] = "user"
	}

	claims["exp"] = time.Now().Add(time.Minute).Unix() // Set expiration time to 1 minute
	// claims["exp"] = time.Now().AddDate(0, 0, 30).Unix() // Set expiration time to 30 days
	// claims["exp"] = time.Now().AddDate(5, 0, 0).Unix() // Set expiration time to 5 years

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign the token"})
		return
	}

	role := "user"
	if storeRole == "3" {
		role = "admin"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Login successful",
		"data": gin.H{
			"username":     user.UserLogin,
			"access_token": tokenString,
			"role":         role, // Assigns 'admin' for role 3, otherwise 'user'
		},
	})
	fmt.Println(user.UserLogin, " Success Login")
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

	hash, _ := HashPassword(user.UserPassword) // ignore error for the sake of simplicity

	fmt.Println("Password:", user.UserPassword)
	fmt.Println("Hash:    ", hash)

	// Insert user to the database
	_, err = db.Exec("INSERT INTO user (user_login, user_password, user_citizen_id, user_role) VALUES (?, ?, ?, ?)",
		user.UserLogin, hash, user.UserCitizenId, user.UserRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success"})
}
