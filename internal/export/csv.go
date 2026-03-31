package export

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"companyEsgDb/internal/entities"
)

func WriteCompaniesCSV(w http.ResponseWriter, companies []entities.Company) error {
	filename := "companies_" + time.Now().Format("20060102_150405") + ".csv"

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	header := []string{
		"id", "name", "bin", "website", "email", "city", "number", "address", "industry", "status",
		"category", "procurement_method", "procurement_email", "procurement_phone",
		"hr_name", "hr_email", "hr_phone",
		"esg_name", "esg_email", "esg_phone", "esg_report_url", "has_esg_dept",
		"linkedin", "facebook", "last_source", "updated_at",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, c := range companies {
		categoryName := ""
		if c.Category != nil {
			categoryName = c.Category.Name
		}

		record := []string{
			strconv.FormatInt(c.ID, 10),
			c.Name,
			c.BIN,
			c.Website,
			c.Email,
			c.City,
			c.Number,
			c.Address,
			c.Industry,
			c.Status,
			categoryName,
			c.ProcurementMethod,
			c.ProcurementEmail,
			c.ProcurementPhone,
			c.HRName,
			c.HREmail,
			c.HRPhone,
			c.ESGName,
			c.ESGEmail,
			c.ESGPhone,
			c.ESGReportURL,
			strconv.FormatBool(c.HasESGDept),
			c.Linkedin,
			c.Facebook,
			c.LastSource,
			c.UpdatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
