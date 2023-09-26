package dao

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbName = "hashbash"

func AddDatabaseFlags(flags *pflag.FlagSet) {
	flags.String("database-host", "localhost", "The hostname or IP address of the hashbash database")
	flags.String("database-ssl-mode", "require", "Whether to use SSL in connecting to the database")
	flags.String("database-username", "postgres", "The username with which to authenticate to the database")
	flags.String("database-password", "postgres", "The password with which to authenticate to the database")
}

func GetConnection() (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database-host"),
		viper.GetString("database-username"),
		viper.GetString("database-password"),
		dbName,
		viper.GetString("database-ssl-mode"),
	)

	return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
}

func GetConnectionOrDie() *gorm.DB {
	db, err := GetConnection()
	if err != nil {
		log.Error().Err(err).Msg("Error creating database connection")
		os.Exit(1)
	}

	return db
}
