package services

import (
	"context"
	"time"

	"companyEsgDb/internal/entities"
	"companyEsgDb/internal/parser"
	"companyEsgDb/internal/repositories"
)

type ParserService struct {
	parser      *parser.WebsiteParser
	companyRepo repositories.CompanyRepository
}

func NewParserService(parser *parser.WebsiteParser, companyRepo repositories.CompanyRepository) *ParserService {
	return &ParserService{parser: parser, companyRepo: companyRepo}
}

func (s *ParserService) ParseAndEnrich(company *entities.Company) error {
	if company.Website == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := s.parser.Parse(ctx, company.Website)
	if err != nil {
		return err
	}

	oldEmail := company.Email
	oldNumber := company.Number

	if company.Email == "" && len(result.Emails) > 0 {
		company.Email = result.Emails[0]
	}
	if company.Number == "" && len(result.Phones) > 0 {
		company.Number = result.Phones[0]
	}
	if company.BIN == "" {
		company.BIN = result.BIN
	}
	if company.Linkedin == "" {
		company.Linkedin = result.Linkedin
	}
	if company.Facebook == "" {
		company.Facebook = result.Facebook
	}

	company.ProcurementMethod = result.ProcurementMethod
	company.ProcurementEmail = result.ProcurementEmail
	company.ProcurementPhone = result.ProcurementPhone
	company.HREmail = result.HREmail
	company.HRPhone = result.HRPhone
	company.ESGEmail = result.ESGEmail
	company.ESGPhone = result.ESGPhone
	company.ESGReportURL = result.ESGReportURL
	company.HasESGDept = result.HasESGDept
	company.LastSource = result.Source
	now := time.Now()
	company.LastParsedAt = &now

	if err := s.companyRepo.UpdateCompany(company); err != nil {
		return err
	}

	if oldEmail != company.Email {
		_ = s.companyRepo.AddLog(&entities.CompanyLog{
			CompanyID: company.ID,
			Action:    "update",
			FieldName: "email",
			OldValue:  oldEmail,
			NewValue:  company.Email,
			Source:    result.Source,
		})
	}
	if oldNumber != company.Number {
		_ = s.companyRepo.AddLog(&entities.CompanyLog{
			CompanyID: company.ID,
			Action:    "update",
			FieldName: "number",
			OldValue:  oldNumber,
			NewValue:  company.Number,
			Source:    result.Source,
		})
	}

	return nil
}
