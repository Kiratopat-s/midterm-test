package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(uid uint, username string, firstName string, lastName string,position string,photoLink string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"firstName": firstName,
		"lastName": lastName,
		"position": position,
		"photoLink": photoLink,
		"exp":      time.Now().Add(time.Minute*30).Unix(),
	})

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing key")
		return signedToken, err
	}
	fmt.Print("Token || ", signedToken)
	return signedToken, nil
}


