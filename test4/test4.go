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
	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.Default.LogMode(logger.Silent)

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

	userCount := viper.GetInt("user-count")
	batchSize := viper.GetInt("batch-size")

	slog.Debug("params", "batchSize", batchSize, "userCount", userCount)

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

	// Query using GORM
	var results []struct {
		Username string
		Count    int
	}

	db.Model(&User{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Having("COUNT(*) > 1").
		Scan(&results)

	// Print or process the results
	for _, result := range results {
		fmt.Printf("Field1: %s,Count: %d\n", result.Username, result.Count)
	}

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
