package common

import (
	"fmt"

	"golang.org/x/text/language"
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

func ShowStats(db *gorm.DB, tableName string) error {
	var rowCount int64
	var fileSize int64

	rowCount, err := getRowCount(db, tableName)
	if err != nil {
		return fmt.Errorf("error getting row count: %w", err)
	}

	// Get file size
	if err := db.Raw("PRAGMA page_size").Scan(&fileSize).Error; err != nil {
		return fmt.Errorf("error getting file size: %w", err)
	}

	p := message.NewPrinter(language.English)
	formattedRowCount := p.Sprintf("%d", rowCount)
	formattedFileSize := p.Sprintf("%d bytes", fileSize)

	fmt.Printf("Row count: %s\n", formattedRowCount)
	fmt.Printf("File size: %s\n", formattedFileSize)

	return nil
}
