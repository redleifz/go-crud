package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func VerifyTokenMiddleware() gin.HandlerFunc {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// fmt.Println(`token :`, tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// fmt.Println("Incoming token:", tokenString) // Log incoming token

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			fmt.Println("Token validation error:", err) // Log token validation error
			// fmt.Println("Token content:", tokenString)  // Log token content

			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			// c.Abort()
			// return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Set the claims in the context for role verification
		c.Set("jwtClaims", claims)
		c.Next()
	}
}

func CheckUserRole(c *gin.Context) {
	claims, _ := c.Get("jwtClaims") // Assuming "jwtClaims" is set in VerifyTokenMiddleware

	//print "claims

	fmt.Println(claims)

	// Extract role from claims (adjust according to your token structure)
	if role, ok := claims.(jwt.MapClaims)["role"].(string); ok {
		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not authorized"})
			c.Abort()
			return
		}
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found or not an integer"})
		c.Abort()
		return
	}
}
