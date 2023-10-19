// controllers/user_controller.go

package controllers

import (
	"go-crud/configs"
	"go-crud/models"
	"go-crud/responses"
	"net/http"

	"github.com/gin-gonic/gin"
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

//PostUser

func CreateUser(c *gin.Context) {
	db := configs.ConnectDB()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert user to database
	_, err := db.Exec("INSERT INTO user (user_login, user_pwd ,user_IDcard) VALUES (?, ?,?)", user.UserLogin, user.UserPwd, user.UserIdCard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success"})
}
