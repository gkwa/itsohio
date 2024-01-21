package test2

import (
	"fmt"
	"log/slog"

	"github.com/taylormonacelli/itsohio/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
}

var (
	UserCount = 1_000
	BatchSize = 1_000
)

func Test2() error {
	db, err := gorm.Open(sqlite.Open("test2.sqlite"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("error auto-migrating User model: %w", err)
	}

	var users []User

	for i := 1; i <= UserCount; i++ {
		username := fmt.Sprintf("user%d", i)
		users = append(users, User{Username: username})

		if i%BatchSize == 0 || i == 50_000 {
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
		TableName: "users",
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
