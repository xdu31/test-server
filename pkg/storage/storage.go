package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	db     *gorm.DB
	onceDB sync.Once
)

// DB returns the current connection instance
func DB() *gorm.DB {
	return db
}

// DBConnect is a more-than-once version of the function that connects the
// executing process with the DB server; directly usable mainly in unit
// tests; indirectly used with the once-wrapper below.
func DBConnect() {
	if viper.GetString("database.dsn") == "" {
		setDBConnection(
			viper.GetString("database.address"),
			viper.GetString("database.port"),
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.ssl"),
			viper.GetString("database.name"),
		)
	}

	var err error

	for i := 0; i < viper.GetInt("database.retries"); i++ {
		log.Debugf("Connecting to database service %s:%s, attempt #%d...",
			viper.GetString("database.address"), viper.GetString("database.port"), i+1)

		db, err = gorm.Open(viper.GetString("database.type"), viper.GetString("database.dsn"))
		if err == nil {
			log.Debugf("Successfully connected to database service %s:%s",
				viper.GetString("database.address"), viper.GetString("database.port"))
			break
		}

		log.Errorf("Database Connection Error: %v", err)
		time.Sleep(time.Duration(time.Duration(i+1) * time.Second))
	}

	if err != nil {
		panic(err)
	}

	db.LogMode(viper.GetBool("database.log"))
}

// Open creates a DB connection
func Open() {
	onceDB.Do(func() {
		DBConnect()
	})
}

// Close closes the DB connection.
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
	log.Info("DB closed")
}

// WithTx performs a function in a separate DB transaction
func WithTx(f func(db *gorm.DB) error) error {
	tx := DB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := f(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			logrus.Errorf("Failed to rollback DB transaction: %v", err)
		}

		return err
	}

	return tx.Commit().Error
}

func setDBConnection(host, port, user, passw, sslmode, dbname string) {
	viper.Set(
		"database.dsn",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s sslmode=%s dbname=%s",
			host, port, user, passw, sslmode, dbname,
		),
	)
}

// OpenFromEnv creates DB connection by environment args for using in tests.
func OpenFromEnv() {
	setDBConnection(
		viper.GetString("PGHOST"),
		viper.GetString("PGPORT"),
		viper.GetString("PGUSER"),
		viper.GetString("PGPASSWORD"),
		"disable",
		viper.GetString("PGDATABASE"),
	)

	viper.Set("database.retries", 3)
	viper.Set("database.log", true)
	viper.Set("database.type", "postgres")

	log.Debug("DB conn: ", viper.GetString("database.dsn"))

	Open()
}
