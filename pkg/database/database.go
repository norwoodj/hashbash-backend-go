package database

import (
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const dbName = "hashbash"

type PageConfig struct {
	Descending bool
	PageNumber int
	PageSize   int
	SortKey    string
}

func AddDatabaseFlags(flags *pflag.FlagSet) {
	flags.StringP("database-host", "d", "", "The hostname or IP address of the hashbash database")
	flags.StringP("database-username", "u", "", "The username with which to authenticate to the database")
	flags.StringP("database-password", "p", "", "The password with which to authenticate to the database")
}

func GetConnection() (*gorm.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = viper.GetString("database-username")
	cfg.Passwd = viper.GetString("database-password")
	cfg.Addr = viper.GetString("database-host")
	cfg.DBName = dbName
	cfg.ParseTime = true

	databaseDsn := cfg.FormatDSN()
	db, err := gorm.Open("mysql", databaseDsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetConnectionOrDie() *gorm.DB {
	db, err := GetConnection()
	if err != nil {
		log.Errorf("Error creating database connection: %s", err)
		os.Exit(1)
	}

	return db
}

func ApplyPaging(db *gorm.DB, pageConfig PageConfig) *gorm.DB {
	orderClause := pageConfig.SortKey
	if pageConfig.Descending {
		orderClause += " DESC"
	}

	return db.Limit(pageConfig.PageSize).
		Offset(pageConfig.PageNumber * pageConfig.PageSize).
		Order(orderClause)
}
