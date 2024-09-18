package user

import (
	"errors"
	"fmt"

	"github.com/Kiratopat-s/workflow/internal/auth"
	"github.com/Kiratopat-s/workflow/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)



type Service struct {
	Repository Repository
	secret     string
}

func NewService(db *gorm.DB, secret string) Service {
	return Service{
		Repository: NewRepository(db),
		secret:     secret,
	}
}

func (service Service) Login(req model.RequestLogin) (string, error) {
	
	user, err := service.Repository.FindOneByUsername(req.Username)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid password")
	}
	token, err := auth.CreateToken(user.ID,user.Username,user.FirstName,user.LastName,user.Position,user.PhotoLink)
	
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service Service) Register(req model.RequestRegister) error {
	fmt.Print("SVR || ",req.Username)
	// Check if username is already taken
	_, err := service.Repository.FindOneByUsername(req.Username)
	if err != nil {
		return errors.New("username already taken")
	}

	// Hash password
	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	// Create user from type
	user := model.User{
		Username: req.Username,
		Password: hash,
		Position: req.Position,
		FirstName: req.FirstName,
		LastName: req.LastName,
		PhotoLink: req.PhotoLink,
	}
	if err := service.Repository.Register(&user); err != nil {
		return err
	}
	return nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Print(err)
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

