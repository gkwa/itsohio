package common

import (
	"fmt"
	"html/template"
	"os"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

func getRowCount(db *gorm.DB, tableName string) (int64, error) {
	var rowCount int64
	var err error

	if db.Migrator().HasTable(tableName) {
		// If the table exists, count rows directly
		err = db.Table(tableName).Count(&rowCount).Error
	} else {
		// If the table doesn't exist, assume it's a model and use Model function
		err = db.Model(tableName).Count(&rowCount).Error
	}

	if err != nil {
		return 0, fmt.Errorf("error getting row count for table %s: %v", tableName, err)
	}

	return rowCount, nil
}

// Define a template for the output
const statsTemplate = `
Row count: {{ .RowCount | formatNumber }}
File size: {{ .FileSize | formatBytes }}
`

// StatsData represents the data for the template
type StatsData struct {
	RowCount   int64
	FileSize   int64
	TableName  string
	DbFilePath string
}

// formatNumber formats the number with thousand separators
func formatNumber(n int64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprint(n)
}

// formatNumber formats the number with thousand separators
func formatBytes(n int64) string {
	bytes := int64(n)
	size := humanize.Bytes(uint64(bytes))
	return size
}

func ShowStats(db *gorm.DB, stats StatsData) error {
	// Get row count
	rc, err := getRowCount(db, stats.TableName)
	if err != nil {
		return fmt.Errorf("error getting row count: %w", err)
	}

	stats.RowCount = rc

	// Get file size
	if err := db.Raw("PRAGMA page_size").Scan(&stats.FileSize).Error; err != nil {
		return fmt.Errorf("error getting file size: %w", err)
	}

	// Get file information
	fileInfo, err := os.Stat(stats.DbFilePath)
	if err != nil {
		return fmt.Errorf("error getting file information: %w", err)
	}

	// Get file size in bytes
	stats.FileSize = fileInfo.Size()

	// Create a new template and parse it
	funcMap := template.FuncMap{
		"formatNumber": formatNumber,
		"formatBytes":  formatBytes,
	}

	tmpl, err := template.New("stats").Funcs(funcMap).Parse(statsTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Execute the template
	if err := tmpl.Execute(os.Stdout, stats); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
