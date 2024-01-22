package test5

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"github.com/taylormonacelli/bravelock/filename"
	"github.com/taylormonacelli/itsohio/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
}

func Test5() error {
	userCount := viper.GetInt("user-count")
	batchSize := viper.GetInt("batch-size")
	gormLogLevel := parseLogLevel(viper.GetString("gorm-log-level"))

	slog.Debug("log level", "gormLogLevel", gormLogLevel)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

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

	for i := 1; i <= userCount; i++ {
		username := fmt.Sprintf("user%d", i)
		users = append(users, User{Username: username})
	}

	db.Create(&users)

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

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		fmt.Println("Invalid log level. Supported values: silent, warn, error, info")
		os.Exit(1)
		return logger.Info
	}
}
