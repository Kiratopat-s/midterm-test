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
	// TABLE users (
	// 	id          SERIAL PRIMARY KEY,
	// 	username    VARCHAR(255) UNIQUE NOT NULL,
	// 	password    VARCHAR(255) NOT NULL,
	// 	position    VARCHAR(100),
	// 	first_name  VARCHAR(100),
	// 	last_name   VARCHAR(100),
	// 	photo_link  TEXT)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      uid,
		"username": username,
		"firstName": firstName,
		"lastName": lastName,
		"position": position,
		"photoLink": photoLink,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing key")
		return signedToken, err
	}
	fmt.Print("Token || ", signedToken)
	return signedToken, nil
}

// func CreateToken(uid  ,username string, secret string) (string, error) {
// 	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
// 		Audience:  jwt.ClaimStrings{uid,username},
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
// 	})

// 	signedToken, err := t.SignedString([]byte(secret))
// 	if err != nil {
// 		log.Println("error signing key")
// 		return signedToken, err
// 	}
// 	fmt.Print("Token || ",signedToken)
// 	return signedToken, nil
// }


// สอบถามเพิ่มเติมส่วนของ file structure ครับ : ส่วนไหนจะเป็น folder ที่ next จะเข้าใจว่าเป็น server side พวกไฟล์ *.ts ที่เราต้องการเขียนเป็น API ครับ 
//npm yarn pnpm แนะนำอันไหนดีครับ