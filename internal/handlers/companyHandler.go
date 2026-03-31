package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"companyEsgDb/internal/entities"
	"companyEsgDb/internal/export"
	"companyEsgDb/internal/repositories"
	"companyEsgDb/internal/services"

	"github.com/gorilla/mux"
)

type CompanyHandler struct {
	companyService *services.CompanyService
	categoryRepo   repositories.CategoryRepository
	templates      *template.Template
}

func NewCompanyHandler(companyService *services.CompanyService, categoryRepo repositories.CategoryRepository) *CompanyHandler {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	return &CompanyHandler{
		companyService: companyService,
		categoryRepo:   categoryRepo,
		templates:      tmpl,
	}
}

func (h *CompanyHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", h.ListCompanies).Methods(http.MethodGet)
	r.HandleFunc("/companies/new", h.ShowCreateForm).Methods(http.MethodGet)
	r.HandleFunc("/companies", h.CreateCompany).Methods(http.MethodPost)
	r.HandleFunc("/companies/import", h.ShowImportForm).Methods(http.MethodGet)
	r.HandleFunc("/companies/import", h.ImportCSV).Methods(http.MethodPost)
	r.HandleFunc("/companies/export/csv", h.ExportCSV).Methods(http.MethodGet)
	r.HandleFunc("/companies/export/excel", h.ExportExcel).Methods(http.MethodGet)
	r.HandleFunc("/companies/{id:[0-9]+}/parse", h.ParseCompany).Methods(http.MethodPost)
	r.HandleFunc("/companies/{id:[0-9]+}/delete", h.DeleteCompany).Methods(http.MethodPost)
	r.HandleFunc("/companies/{id}", h.ViewCompany).Methods(http.MethodGet)
	r.HandleFunc("/companies/delete-all", h.DeleteAllCompanies).Methods(http.MethodPost)
}

func (h *CompanyHandler) buildFilter(r *http.Request) repositories.CompanyFilter {
	var hasESG *bool
	if r.URL.Query().Get("has_esg") == "true" {
		v := true
		hasESG = &v
	}

	categoryID, _ := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
	return repositories.CompanyFilter{
		Search:            r.URL.Query().Get("search"),
		City:              r.URL.Query().Get("city"),
		CategoryID:        categoryID,
		ProcurementMethod: r.URL.Query().Get("procurement_method"),
		HasESGDept:        hasESG,
	}
}

func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.categoryRepo.GetAll()
	filter := h.buildFilter(r)
	companies, _ := h.companyService.List(filter)

	data := map[string]any{
		"Companies":  companies,
		"Categories": categories,
	}
	_ = h.templates.ExecuteTemplate(w, "index.html", data)
}

func (h *CompanyHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.categoryRepo.GetAll()
	data := map[string]any{"Categories": categories}
	_ = h.templates.ExecuteTemplate(w, "create.html", data)
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	categoryID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)

	company := &entities.Company{
		Name:              r.FormValue("name"),
		BIN:               r.FormValue("bin"),
		Website:           r.FormValue("website"),
		Email:             r.FormValue("email"),
		City:              r.FormValue("city"),
		Number:            r.FormValue("number"),
		Address:           r.FormValue("address"),
		CategoryID:        categoryID,
		ProcurementMethod: r.FormValue("procurement_method"),
		ProcurementEmail:  r.FormValue("procurement_email"),
		ProcurementPhone:  r.FormValue("procurement_phone"),
		HRName:            r.FormValue("hr_name"),
		HREmail:           r.FormValue("hr_email"),
		HRPhone:           r.FormValue("hr_phone"),
		ESGName:           r.FormValue("esg_name"),
		ESGEmail:          r.FormValue("esg_email"),
		ESGPhone:          r.FormValue("esg_phone"),
		ESGReportURL:      r.FormValue("esg_report_url"),
		HasESGDept:        r.FormValue("has_esg") == "on",
	}

	_ = h.companyService.Create(company)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CompanyHandler) ShowImportForm(w http.ResponseWriter, r *http.Request) {
	_ = h.templates.ExecuteTemplate(w, "import.html", nil)
}

func (h *CompanyHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	path, err := h.companyService.SaveUploadedFile(file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.companyService.ImportCSV(path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CompanyHandler) ParseCompany(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	_ = h.companyService.ParseCompany(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	_ = h.companyService.Delete(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CompanyHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	filter := h.buildFilter(r)
	companies, err := h.companyService.List(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := export.WriteCompaniesCSV(w, companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CompanyHandler) ExportExcel(w http.ResponseWriter, r *http.Request) {
	filter := h.buildFilter(r)
	companies, err := h.companyService.List(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := export.WriteCompaniesExcel(w, companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CompanyHandler) ViewCompany(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	company, err := h.companyService.GetByID(id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	err = h.templates.ExecuteTemplate(w, "view.html", company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CompanyHandler) DeleteAllCompanies(w http.ResponseWriter, r *http.Request) {
	err := h.companyService.DeleteAllCompanies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
