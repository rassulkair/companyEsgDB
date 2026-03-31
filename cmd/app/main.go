package main

import (
	"companyEsgDb/internal/config"
	"companyEsgDb/internal/db"
	"companyEsgDb/internal/handlers"
	"companyEsgDb/internal/jobs"
	"companyEsgDb/internal/migrations"
	"companyEsgDb/internal/parser"
	"companyEsgDb/internal/repositories"
	"companyEsgDb/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	if err := migrations.Run(cfg.PostgresURL()); err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	companyRepo := repositories.NewCompanyRepository(database)
	categoryRepo := repositories.NewCategoryRepository(database)
	websiteParser := parser.NewWebsiteParser()
	parserService := services.NewParserService(websiteParser, companyRepo)
	companyService := services.NewCompanyService(companyRepo, categoryRepo, parserService)
	handler := handlers.NewCompanyHandler(companyService, categoryRepo)

	autoRefreshJob := jobs.NewAutoRefreshJob(companyService, 12*time.Hour)
	autoRefreshJob.Start()

	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))),
	)

	log.Printf("server started on :%s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, router); err != nil {
		log.Fatal(err)
	}
}
