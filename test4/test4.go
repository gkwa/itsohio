package test4

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/taylormonacelli/bravelock/filename"
	"github.com/taylormonacelli/itsohio/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
}

func Test4() error {
	userCount := viper.GetInt("user-count")
	batchSize := viper.GetInt("batch-size")

	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.Default.LogMode(logger.Info)

	strategy := &filename.FilenameFromGoPackageStrategy{}
	fname := filename.GetFnameWithoutExtension(strategy.GetFilename())
	fname = fmt.Sprintf("%s.sqlite", fname)

	dialector := sqlite.Open(fname)
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("error auto-migrating User model: %w", err)
	}

	var users []User

	slog.Debug("params", "batchSize", batchSize, "userCount", userCount)

	var results []struct {
		Username string
		Count    int
	}

	// Reset or zero out the results slice
	results = []struct {
		Username string
		Count    int
	}{}

	db.Model(&User{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Having("COUNT(*) > 1").
		Scan(&results)

	// Print or process the results
	for _, result := range results {
		fmt.Printf("Username: %s, Count: %d\n", result.Username, result.Count)
	}

	// Reset or zero out the results slice
	results = []struct {
		Username string
		Count    int
	}{}

	for i := 1; i <= userCount; i++ {
		username := fmt.Sprintf("user%d", i)
		users = append(users, User{Username: username})

		if i%batchSize == 0 || i == userCount {
			// Insert the batch
			result := db.Create(&users)
			if result.Error != nil {
				slog.Error("error inserting users", "error", result.Error)
				return fmt.Errorf("error inserting users: %v", result.Error)
			}

			// Clear the slice for the next batch
			users = []User{}
		}
	}

	db.Model(&User{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Having("COUNT(*) > 1").
		Scan(&results)

	// Print or process the results
	for _, result := range results {
		fmt.Printf("Username: %s, Count: %d\n", result.Username, result.Count)
	}

	// Reset or zero out the results slice
	results = []struct {
		Username string
		Count    int
	}{}

	// Query duplicates using GORM
	var duplicateUsernames []string
	db.Model(&User{}).
		Select("username").
		Group("username").
		Having("COUNT(*) > 1").
		Pluck("username", &duplicateUsernames)

	slog.Debug("duplicateUsernames", "count", len(duplicateUsernames))

	// Create a map to store IDs of records to be deleted
	recordsToDelete := make(map[uint]bool)

	// Loop through each duplicate username
	for _, username := range duplicateUsernames {
		// Find the records with the duplicate username and order by created_at ascending
		var duplicateRecords []User
		db.Where("username = ?", username).
			Order("created_at ASC").
			Find(&duplicateRecords)

		// Keep the first record (oldest created_at), add others to the deletion list
		for i := 1; i < len(duplicateRecords); i++ {
			recordsToDelete[duplicateRecords[i].ID] = true
		}
	}

	// Convert map keys (record IDs) to a slice
	var idsToDelete []uint
	for id := range recordsToDelete {
		idsToDelete = append(idsToDelete, id)
	}

	// Delete records in batches
	batchSizeDelete := 100
	for i := 0; i < len(idsToDelete); i += batchSizeDelete {
		end := i + batchSizeDelete
		if end > len(idsToDelete) {
			end = len(idsToDelete)
		}

		batchIDs := idsToDelete[i:end]

		// Perform a single delete statement for each batch
		result := db.Where("id IN ?", batchIDs).Delete(&User{})
		if result.Error != nil {
			slog.Error("error deleting duplicate users", "error", result.Error)
			return fmt.Errorf("error deleting duplicate users: %v", result.Error)
		}
	}

	// Reset or zero out the results slice
	results = []struct {
		Username string
		Count    int
	}{}

	db.Model(&User{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Scan(&results)

	// Print or process the results
	for _, result := range results {
		fmt.Printf("Username: %s, Count: %d\n", result.Username, result.Count)
	}

	// Reset or zero out the results slice
	results = []struct {
		Username string
		Count    int
	}{}

	stats := common.StatsData{
		TableName:  "users",
		DbFilePath: fname,
	}

	err = common.ShowStats(db, stats)
	if err != nil {
		return fmt.Errorf("error showing stats: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error getting database connection: %w", err)
	}
	sqlDB.Close()

	return nil
}
