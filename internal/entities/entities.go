package entities

import "time"

type Company struct {
	ID int64 `gorm:"primaryKey"`

	Name     string `gorm:"size:255;not null"`
	BIN      string `gorm:"size:20;index"`
	Website  string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	City     string `gorm:"size:120"`
	Number   string `gorm:"size:120"`
	Address  string `gorm:"type:text"`
	Industry string `gorm:"size:255"`
	Status   string `gorm:"size:50;default:active"`

	DirectorName  string `gorm:"size:255"`
	DirectorPos   string `gorm:"size:255"`
	ExecutiveName string `gorm:"size:255"`
	ExecutivePos  string `gorm:"size:255"`
	DirStart      string `gorm:"size:100"`
	ExecStart     string `gorm:"size:100"`

	Linkedin     string `gorm:"size:255"`
	Facebook     string `gorm:"size:255"`
	StatusFb     string `gorm:"size:100"`
	StatusLink   string `gorm:"size:100"`
	LiLastUpdate string `gorm:"size:100"`
	FbLastUpdate string `gorm:"size:100"`

	ProcurementMethod string `gorm:"size:255"`
	ProcurementEmail  string `gorm:"size:255"`
	ProcurementPhone  string `gorm:"size:120"`

	HRName  string `gorm:"size:255"`
	HREmail string `gorm:"size:255"`
	HRPhone string `gorm:"size:120"`

	ESGName      string `gorm:"size:255"`
	ESGEmail     string `gorm:"size:255"`
	ESGPhone     string `gorm:"size:120"`
	ESGReportURL string `gorm:"size:255"`
	HasESGDept   bool   `gorm:"default:false"`
	LastSource   string `gorm:"size:255"`

	LastParsedAt *time.Time
	UpdatedAt    time.Time
	CreatedAt    time.Time

	CategoryID int64        `gorm:"index"`
	Category   *Category    `gorm:"foreignKey:CategoryID"`
	Logs       []CompanyLog `gorm:"foreignKey:CompanyID"`
}

type Category struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:255;uniqueIndex;not null"`
}

type CompanyLog struct {
	ID        int64 `gorm:"primaryKey"`
	CompanyID int64 `gorm:"index;not null"`

	Action    string `gorm:"size:100"`
	FieldName string `gorm:"size:100"`
	OldValue  string `gorm:"type:text"`
	NewValue  string `gorm:"type:text"`
	Source    string `gorm:"size:255"`
	CreatedAt time.Time

	Company *Company `gorm:"foreignKey:CompanyID"`
}
