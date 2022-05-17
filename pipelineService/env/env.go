package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type envFile struct {
	BuildEnv                 string
	ServerPort               string
	AirByteAddress           string
	AuthServiceAddress       string
	DbHost                   string
	DbUsername               string
	DbPassword               string
	DbName                   string
	DbPort                   string
	CadenceService           string
	TaskListName             string
	CadenceDomain            string
	ClientName               string
	CadenceServiceName       string
	CadenceWorkerServiceName string
	DefaultCSVSourcePath     string
}

var Env *envFile

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}

	buildEnv := os.Getenv("BUILD_ENV")
	if buildEnv == "" {
		buildEnv = "dev"
	}

	airByteAddress := os.Getenv("AIR_BYTE_ADDRESS")
	if airByteAddress == "" {
		airByteAddress = "http://localhost:8000"
	}

	authServiceAddress := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authServiceAddress == "" {
		authServiceAddress = "http://localhost:8002"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		dbUsername = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	defaultCSVSourcePath := os.Getenv("DEFAULT_CSV_SOURCE_PATH")
	if defaultCSVSourcePath == "" {
		defaultCSVSourcePath = "/local/temp.csv"
	}

	Env = &envFile{
		BuildEnv:                 buildEnv,
		ServerPort:               serverPort,
		AirByteAddress:           airByteAddress,
		AuthServiceAddress:       authServiceAddress,
		DbHost:                   dbHost,
		DbUsername:               dbUsername,
		DbPassword:               dbPassword,
		DbName:                   dbName,
		DbPort:                   dbPort,
		CadenceService:           os.Getenv("CADENCE_SERVICE"),
		TaskListName:             os.Getenv("TASK_LIST_NAME"),
		CadenceDomain:            os.Getenv("CADENCE_DOMAIN"),
		ClientName:               os.Getenv("CLIENT_NAME"),
		CadenceServiceName:       os.Getenv("CADENCE_SERVICE_NAME"),
		CadenceWorkerServiceName: os.Getenv("CADENCE_WORKER_SERVICE_NAME"),
		DefaultCSVSourcePath:     defaultCSVSourcePath,
	}
}
