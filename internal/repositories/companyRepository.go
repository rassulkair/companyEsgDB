package repositories

import (
	"companyEsgDb/internal/entities"

	"gorm.io/gorm"
)

type CompanyFilter struct {
	Search            string
	CategoryID        int64
	City              string
	ProcurementMethod string
	HasESGDept        *bool
}

type CompanyRepository interface {
	AddCompany(company *entities.Company) error
	GetAllCompanies(filter CompanyFilter) ([]entities.Company, error)
	GetCompanyByID(id int64) (*entities.Company, error)
	GetCompanyByBIN(bin string) (*entities.Company, error)
	UpdateCompany(company *entities.Company) error
	DeleteCompany(id int64) error
	SaveBatch(companies []entities.Company) error
	AddLog(log *entities.CompanyLog) error
	DeleteAllCompanies() error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) AddCompany(company *entities.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) GetAllCompanies(filter CompanyFilter) ([]entities.Company, error) {
	var companies []entities.Company
	query := r.db.Preload("Category").Order("updated_at desc")

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR bin ILIKE ? OR website ILIKE ?", like, like, like)
	}
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.City != "" {
		query = query.Where("city ILIKE ?", "%"+filter.City+"%")
	}
	if filter.ProcurementMethod != "" {
		query = query.Where("procurement_method ILIKE ?", "%"+filter.ProcurementMethod+"%")
	}
	if filter.HasESGDept != nil {
		query = query.Where("has_esg_dept = ?", *filter.HasESGDept)
	}

	err := query.Find(&companies).Error
	return companies, err
}

func (r *companyRepository) GetCompanyByID(id int64) (*entities.Company, error) {
	var company entities.Company
	err := r.db.Preload("Category").First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}
func (r *companyRepository) GetCompanyByBIN(bin string) (*entities.Company, error) {
	var company entities.Company
	err := r.db.Where("bin = ?", bin).First(&company).Error
	return &company, err
}

func (r *companyRepository) UpdateCompany(company *entities.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) DeleteCompany(id int64) error {
	return r.db.Delete(&entities.Company{}, id).Error
}

func (r *companyRepository) SaveBatch(companies []entities.Company) error {
	return r.db.CreateInBatches(companies, 100).Error
}

func (r *companyRepository) AddLog(logEntry *entities.CompanyLog) error {
	return r.db.Create(logEntry).Error
}

func (r *companyRepository) DeleteAllCompanies() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entities.Company{}).Error
}
