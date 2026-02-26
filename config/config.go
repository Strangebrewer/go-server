package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBUsername    string `json:"mongo_username"`
	MongoDBPassword    string `json:"mongo_password"`
	MongoDBCluster     string `json:"mongo_cluster"`
	MongoDBName        string `json:"mongo_db_name"`
	Database           string `json:"database"`
	DBUser             string `json:"db_user"`
	DBPassword         string `json:"db_password"`
	DBName             string `json:"db_name"`
	DBHost             string `json:"db_host"`
	DBPort             string `json:"db_port"`
	PrivateKeyPEM      string
	PublicKeyPEM       string
	RefreshTokenPepper string
	loaded             bool
}

var Cfg Config

func InitConfig() *Config {

	cfg := &Config{
		loaded: false,
	}

	filename := "./.env"
	if _, err := os.Stat(filename); err == nil {
		_ = godotenv.Load(filename)
	}

	return cfg
}

func GetCurrentCfg() Config {
	return Cfg
}

func (cfg *Config) LoadEnvVariables() {
	cfg.PrivateKeyPEM = os.Getenv("PRIVATE_KEY")
	cfg.PublicKeyPEM = os.Getenv("PUBLIC_KEY")

	refreshPepper := os.Getenv("REFRESH_TOKEN_PEPPER")
	if refreshPepper != "" {
		cfg.RefreshTokenPepper = refreshPepper
	}

	mongoDbUsername := os.Getenv("MONGO_USERNAME")
	if mongoDbUsername != "" {
		cfg.MongoDBUsername = mongoDbUsername
	}

	mongoDbpassword := os.Getenv("MONGO_PASSWORD")
	if mongoDbpassword != "" {
		cfg.MongoDBPassword = mongoDbpassword
	}

	mongoDbCluster := os.Getenv("MONGO_CLUSTER")
	if mongoDbCluster != "" {
		cfg.MongoDBCluster = mongoDbCluster
	}

	mongoDbName := os.Getenv("MONGO_DB_NAME")
	if mongoDbName != "" {
		cfg.MongoDBName = mongoDbName
	}

	database := os.Getenv("DATABASE")
	if database != "" {
		cfg.Database = database
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser != "" {
		cfg.DBUser = dbUser
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword != "" {
		cfg.DBPassword = dbPassword
	}

	dbName := os.Getenv("DB_NAME")
	if dbName != "" {
		cfg.DBName = dbName
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		cfg.DBHost = dbHost
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort != "" {
		cfg.DBPort = dbPort
	}

	cfg.loaded = true
	Cfg = *cfg
}

func (cfg *Config) IsLoaded() bool {
	return cfg.loaded
}
