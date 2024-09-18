package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Guard(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token "Bearer xxx" from cookie
		auth, err := c.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Remove prefix "Bearer " from auth token
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			log.Println("Token does not have the expected 'Bearer ' prefix")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := verifyToken(tokenString, secret)
		if err != nil {
			log.Printf("Token verification failed: %v\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("Token verified successfully. Claims: %+v\n", claims)
			
			keys := []string{"uid", "username", "firstName", "lastName", "position", "photoLink"}

			for _, key := range keys {
				if value, ok := claims[key].(string); ok {
					c.Set(key, value)
				} else if key == "uid" {
					if value, ok := claims[key].(float64); ok {
						c.Set(key, value)
					}
				}
			}
		} else {
			log.Println("Token claims are invalid or token is not valid")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate that the algorithm matches HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(secret), nil
	})
	

	// Check for verification errors
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %v", err)
	}

	log.Println("Token parsed successfully, checking signature validity")

	// Return the verified token
	return token, nil
}

func GuardAdmin(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token "Bearer xxx" from cookie
		auth, err := c.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Remove prefix "Bearer " from auth token
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			log.Println("Token does not have the expected 'Bearer ' prefix")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := verifyToken(tokenString, secret)
		if err != nil {
			log.Printf("Token verification failed: %v\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Printf("Token verified successfully. Claims: %+v\n", claims)
			
			keys := []string{"uid", "username", "firstName", "lastName", "position", "photoLink"}

			for _, key := range keys {
				if value, ok := claims[key].(string); ok {
					c.Set(key, value)
				} else if key == "uid" {
					if value, ok := claims[key].(float64); ok {
						c.Set(key, value)
					}
				}
			}

			if role, ok := claims["position"].(string); !ok || role != "Admin" {
				log.Println("User is not an admin")
				c.AbortWithStatus(http.StatusForbidden) // 403 Forbidden
				return
			}
		} else {
			log.Println("Token claims are invalid or token is not valid")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}