package user

import (
	"github.com/Kiratopat-s/workflow/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	Database *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		Database: db,
	}
}

func (repo Repository) FindOneByUsername(username string) (model.User, error) {
	var result model.User

	db := repo.Database
	db = db.Where("username = ?", username)

	if err := db.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (repo Repository) Register(user *model.User) error {
	db := repo.Database
	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
