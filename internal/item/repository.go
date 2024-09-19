package item

import (
	"errors"

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

func (repo Repository) Create(item *model.Item) error {
	return repo.Database.Create(item).Error
}

func (repo Repository) Find(query model.RequestFindItem) ([]model.Item, error) {
	var results []model.Item

	db := repo.Database

	if statuses := query.Status; len(statuses) > 0 {
		db = db.Where("status = ?", statuses)
	}

	if item_id := query.ItemID; item_id > 0 {
		db = db.Where("id = ?", item_id)
	}

	if err := db.Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (repo Repository) FindAll() ([]model.Item, error) {
	var results []model.Item
	if err := repo.Database.Order("id desc").Find(&results).Error; err != nil {
		return results, err
	}
	return results, nil
}

func (repo Repository) FindByID(id uint) (model.Item, error) {
	var result model.Item

	// Query the database
	err := repo.Database.First(&result, id).Error
	if err != nil {
		// Return "record not found" error if no record is found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, gorm.ErrRecordNotFound
		}
		// Return other errors, if any
		return result, err
	}

	// Return the found result
	return result, nil
}

func (repo Repository) Replace(item model.Item) error {
	return repo.Database.Model(&item).Updates(item).Error
}

func (repo Repository) Delete(id uint) error {
	return repo.Database.Delete(&model.Item{}, id).Error
}

func (repo Repository) UpdateManyStatus(id []int, status string) error {
	return repo.Database.Model(&model.Item{}).Where("id IN (?) AND status = ?", id, "PENDING").Update("status", status).Error
}

func (repo Repository) DeleteMany(id []int) error {
	return repo.Database.Where("id IN (?)", id).Delete(&model.Item{}).Error
}

func (repo Repository) DeleteManyByUserId(ids []int, ownerID int) error {
	return repo.Database.Where("id IN (?) AND owner_id = ?", ids, ownerID).Delete(&model.Item{}).Error
}

func (repo Repository) CountItemsStatusByUser(ownerID int) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := repo.Database.Model(&model.Item{}).Select("status, count(*) as count").Where("owner_id = ?", ownerID).Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]int)
	for _, r := range results {
		resultMap[r.Status] = r.Count
	}

	return resultMap, nil
}