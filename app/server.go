package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB

	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome " + appConfig.AppName)

	server.initializeDB(dbConfig)

	server.initializeRoutes()
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error

	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, string(dbConfig.DBPort), dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, string(dbConfig.DBPort))
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic("failed connect to db")
	}

	for _, model := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated successfully.")
}

func (server *Server) Run(address string) {
	fmt.Println("Listening to port", address)
	log.Fatal(http.ListenAndServe(address, server.Router))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "GoToko")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "admin")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "")
	dbConfig.DBName = getEnv("DB_NAME", "dbname")
	dbConfig.DBPort = getEnv("DB_PORT", "9000")
	dbConfig.DBDriver = getEnv("DB_PORT", "postgre")

	server.Initialize(appConfig, dbConfig)

	server.Run(":" + appConfig.AppPort)
}