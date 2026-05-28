package product

import (
	"order/api/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	Database *db.Db
}

type ProductRepositoryDeps struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) Create(Product *Product) (*Product, error) {
	result := repo.Database.DB.Create(Product)
	if result.Error != nil {
		return nil, result.Error
	}
	return Product, nil
}

func (repo *ProductRepository) GetByID(id uint) (*Product, error) {
	var Product Product
	result := repo.Database.DB.First(&Product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Product, nil
}

func (repo *ProductRepository) Update(Product *Product) (*Product, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(Product)
	if result.Error != nil {
		return nil, result.Error
	}
	return Product, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
