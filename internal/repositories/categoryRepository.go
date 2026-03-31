package repositories

import (
	"companyEsgDb/internal/entities"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetByNameOrCreate(name string) (*entities.Category, error)
	GetAll() ([]entities.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetByNameOrCreate(name string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.FirstOrCreate(&category, entities.Category{Name: name}).Error
	return &category, err
}

func (r *categoryRepository) GetAll() ([]entities.Category, error) {
	var categories []entities.Category
	err := r.db.Order("name asc").Find(&categories).Error
	return categories, err
}
