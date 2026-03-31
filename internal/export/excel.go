package export

import (
	"fmt"
	"net/http"
	"time"

	"companyEsgDb/internal/entities"

	"github.com/xuri/excelize/v2"
)

func WriteCompaniesExcel(w http.ResponseWriter, companies []entities.Company) error {
	file := excelize.NewFile()
	defer func() {
		_ = file.Close()
	}()

	sheet := "Companies"
	file.SetSheetName("Sheet1", sheet)

	headers := []string{
		"ID", "Name", "BIN", "Website", "Email", "City", "Phone", "Address", "Industry", "Status",
		"Category", "Procurement Method", "Procurement Email", "Procurement Phone",
		"HR Name", "HR Email", "HR Phone",
		"ESG Name", "ESG Email", "ESG Phone", "ESG Report URL", "Has ESG Dept",
		"LinkedIn", "Facebook", "Last Source", "Updated At",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = file.SetCellValue(sheet, cell, header)
	}

	style, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	_ = file.SetCellStyle(sheet, "A1", "Z1", style)

	for rowIndex, c := range companies {
		row := rowIndex + 2
		categoryName := ""
		if c.Category != nil {
			categoryName = c.Category.Name
		}

		values := []any{
			c.ID,
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
			c.HasESGDept,
			c.Linkedin,
			c.Facebook,
			c.LastSource,
			c.UpdatedAt.Format(time.RFC3339),
		}

		for colIndex, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, row)
			_ = file.SetCellValue(sheet, cell, value)
		}
	}

	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		_ = file.SetColWidth(sheet, col, col, 20)
	}

	filename := fmt.Sprintf("companies_%s.xlsx", time.Now().Format("20060102_150405"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	return file.Write(w)
}
