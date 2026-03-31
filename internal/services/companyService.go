package services

import (
	"companyEsgDb/internal/entities"
	"companyEsgDb/internal/repositories"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CompanyService struct {
	companyRepo   repositories.CompanyRepository
	categoryRepo  repositories.CategoryRepository
	parserService *ParserService
}

func NewCompanyService(cRepo repositories.CompanyRepository, catRepo repositories.CategoryRepository, parserService *ParserService) *CompanyService {
	return &CompanyService{
		companyRepo:   cRepo,
		categoryRepo:  catRepo,
		parserService: parserService,
	}
}

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	return "https://" + url
}

func (s *CompanyService) Create(company *entities.Company) error {
	if company.CategoryID == 0 && company.Category != nil && company.Category.Name != "" {
		cat, err := s.categoryRepo.GetByNameOrCreate(company.Category.Name)
		if err != nil {
			return err
		}
		company.CategoryID = cat.ID
	}
	if company.Status == "" {
		company.Status = "active"
	}
	return s.companyRepo.AddCompany(company)
}

func (s *CompanyService) List(filter repositories.CompanyFilter) ([]entities.Company, error) {
	return s.companyRepo.GetAllCompanies(filter)
}

func (s *CompanyService) GetByID(id int64) (*entities.Company, error) {
	return s.companyRepo.GetCompanyByID(id)
}

func (s *CompanyService) Delete(id int64) error {
	return s.companyRepo.DeleteCompany(id)
}

func (s *CompanyService) ParseCompany(id int64) error {
	company, err := s.companyRepo.GetCompanyByID(id)
	if err != nil {
		return err
	}
	return s.parserService.ParseAndEnrich(company)
}

func (s *CompanyService) SaveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	path := filepath.Join("uploads", header.Filename)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return path, err
}

func (s *CompanyService) ImportCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	_, _ = reader.Read()

	var companies []entities.Company

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 19 {
			continue
		}

		categoryName := strings.TrimSpace(record[1])
		category, err := s.categoryRepo.GetByNameOrCreate(categoryName)
		if err != nil {
			continue
		}

		hasESG := false
		if len(record) > 21 {
			hasESG, _ = strconv.ParseBool(strings.TrimSpace(record[21]))
		}

		comp := entities.Company{
			Name:       strings.TrimSpace(record[0]),
			CategoryID: category.ID,
			Website:    normalizeURL(strings.TrimSpace(record[2])),
			Email:      strings.TrimSpace(record[3]),
			City:       strings.TrimSpace(record[4]),
			Number:     strings.TrimSpace(record[5]),
			Address:    strings.TrimSpace(record[6]),

			DirectorName:  strings.TrimSpace(record[7]),
			DirectorPos:   strings.TrimSpace(record[8]),
			DirStart:      strings.TrimSpace(record[9]),
			ExecutiveName: strings.TrimSpace(record[10]),
			ExecutivePos:  strings.TrimSpace(record[11]),
			ExecStart:     strings.TrimSpace(record[12]),

			Linkedin:     normalizeURL(strings.TrimSpace(record[13])),
			StatusLink:   strings.TrimSpace(record[14]),
			LiLastUpdate: strings.TrimSpace(record[15]),
			Facebook:     normalizeURL(strings.TrimSpace(record[16])),
			StatusFb:     strings.TrimSpace(record[17]),
			FbLastUpdate: strings.TrimSpace(record[18]),
			Status:       "active",
			HasESGDept:   hasESG,
		}

		companies = append(companies, comp)
	}

	if len(companies) == 0 {
		return fmt.Errorf("no valid companies found in csv")
	}

	return s.companyRepo.SaveBatch(companies)
}

func (s *CompanyService) DeleteAllCompanies() error {
	return s.companyRepo.DeleteAllCompanies()
}
