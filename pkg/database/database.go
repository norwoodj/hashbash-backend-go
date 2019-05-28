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

func AddDatabaseFlags(flags *pflag.FlagSet) {
	flags.String("database-host", "localhost", "The hostname or IP address of the hashbash database")
	flags.String("database-username", "root", "The username with which to authenticate to the database")
	flags.String("database-password", "root", "The password with which to authenticate to the database")
}

func GetConnection() (*gorm.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = viper.GetString("database-username")
	cfg.Passwd = viper.GetString("database-password")
	cfg.Net = "tcp"
	cfg.Addr = viper.GetString("database-host")
	cfg.DBName = dbName
	cfg.ParseTime = true

	databaseDsn := cfg.FormatDSN()
	log.Debug(databaseDsn)
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

func ApplyPaging(
	db *gorm.DB,
	limit int,
	offset int,
	orderColumn string,
	descending bool,
) *gorm.DB {
	orderClause := orderColumn
	if descending {
		orderClause += " DESC"
	}

	return db.Limit(limit).
		Offset(offset).
		Order(orderClause)
}
