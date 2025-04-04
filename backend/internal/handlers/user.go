package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ProgramadoresSemPatria/Team-9/internal/models"
	"github.com/ProgramadoresSemPatria/Team-9/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.ID.String(),
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateUserHandler(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Verified: true,
	}
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}
	c.JSON(http.StatusCreated, models.FilteredResponse(newUser))
}

func LoginHandler(c *gin.Context) {
	var input models.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingUser models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := verifyPassword(existingUser.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := GenerateToken(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		tokenString := c.Request.Header.Get("Authorization")

		if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if _, err := uuid.Parse(claims.Subject); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format in token"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}

func ProfileHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var u models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&u, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, models.FilteredResponse(u))
}
