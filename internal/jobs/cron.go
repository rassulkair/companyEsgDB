package jobs

import (
	"log"
	"time"

	"companyEsgDb/internal/repositories"
	"companyEsgDb/internal/services"
)

type AutoRefreshJob struct {
	companyService *services.CompanyService
	interval       time.Duration
}

func NewAutoRefreshJob(companyService *services.CompanyService, interval time.Duration) *AutoRefreshJob {
	return &AutoRefreshJob{
		companyService: companyService,
		interval:       interval,
	}
}

func (j *AutoRefreshJob) Start() {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		log.Printf("auto-refresh job started, interval=%s", j.interval)

		for range ticker.C {
			j.runOnce()
		}
	}()
}

func (j *AutoRefreshJob) runOnce() {
	companies, err := j.companyService.List(repositories.CompanyFilter{})
	if err != nil {
		log.Printf("auto-refresh list error: %v", err)
		return
	}

	for _, company := range companies {
		if company.Website == "" {
			continue
		}
		if err := j.companyService.ParseCompany(company.ID); err != nil {
			log.Printf("auto-refresh parse error for company id=%d: %v", company.ID, err)
		}
	}

	log.Printf("auto-refresh cycle finished, processed=%d", len(companies))
}
