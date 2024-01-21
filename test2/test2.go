package test2

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"github.com/taylormonacelli/bravelock/filename"
	"github.com/taylormonacelli/itsohio/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
}

func Test2() error {
	strategy := &filename.FilenameFromGoPackageStrategy{}
	fname := filename.GetFnameWithoutExtension(strategy.GetFilename())
	fname = fmt.Sprintf("%s.sqlite", fname)

	db, err := gorm.Open(sqlite.Open(fname), &gorm.Config{})
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
