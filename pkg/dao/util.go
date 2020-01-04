package dao

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const dbName = "hashbash"

// You know, this isn't really that important, and the frontend checks it anyway. so sue me
func isValidHash(hashFunctionName string, hash string) bool {
	return true
}

func AddDatabaseFlags(flags *pflag.FlagSet) {
	flags.String("database-engine", "postgres", "The type of database backend to use (mysql or postgresql)")
	flags.String("database-host", "localhost", "The hostname or IP address of the hashbash database")
	flags.String("database-ssl-mode", "require", "Whether to use SSL in connecting to the database")
	flags.String("database-username", "postgres", "The username with which to authenticate to the database")
	flags.String("database-password", "postgres", "The password with which to authenticate to the database")
}

func getMysqlConnection() (*gorm.DB, error) {
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

func getPostgresqlConnection() (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database-host"),
		viper.GetString("database-username"),
		viper.GetString("database-password"),
		dbName,
		viper.GetString("database-ssl-mode"),
	)

	return gorm.Open("postgres", connectionString)
}

func GetConnection(engine string) (*gorm.DB, error) {
	switch engine {
	case "mysql":
		return getMysqlConnection()
	case "postgres":
		return getPostgresqlConnection()
	case "postgresql":
		return getPostgresqlConnection()
	default:
		return nil, fmt.Errorf("%s engine not supported", engine)
	}
}

func GetConnectionOrDie(engine string) *gorm.DB {
	db, err := GetConnection(engine)
	if err != nil {
		log.Errorf("Error creating database connection: %s", err)
		os.Exit(1)
	}

	return db
}
