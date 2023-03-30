package helpers

import (
	"errors"
	"fmt"
	"job-portal-lite/models/responses"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "key"

func GenerateToken(id int, email string, name string, role string) (string, error) {
	expiry := time.Now().Add(time.Hour * 1).Unix()
	fmt.Println("Expired", expiry)

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"name":  name,
		"role":  role,
		"exp":   expiry,
	}

	// pilih metode enkripsi
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// parsing token menjadi string
	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

//func VerifyToken(c *gin.Context) (interface{}, error) {
//	errResponse := errors.New("sign in to proceed")
//	headerToken := c.Request.Header.Get("Authorization")
//	bearer := strings.HasPrefix(headerToken, "Bearer")
//
//	if !bearer {
//		return nil, errResponse
//	}
//
//	stringToken := strings.Split(headerToken, " ")[1]
//
//	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
//		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errResponse
//		}
//		return []byte(secretKey), nil
//	})
//
//	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
//		return nil, errResponse
//	}
//
//	return token.Claims.(jwt.MapClaims), nil
//}

func VerifyToken(c *gin.Context) (*responses.TokenDecoded, error) {
	errResponse := errors.New("sign in to proceed")
	headerToken := c.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errResponse
	}

	stringToken := strings.Split(headerToken, " ")[1]

	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &responses.TokenDecoded{
			ID:    int(claims["id"].(float64)),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}, nil
	}

	return nil, errResponse
}
