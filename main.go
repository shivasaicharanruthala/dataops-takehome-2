package main

import (
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome-2/database"
	"github.com/shivasaicharanruthala/dataops-takehome-2/etl"
	"github.com/shivasaicharanruthala/dataops-takehome-2/log"
	_ "github.com/shivasaicharanruthala/dataops-takehome-2/model"
	"os"
	"strconv"
	_ "time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// init function runs before the main function. It loads environment variables from a .env file.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	// Retrieve the number of workers and batch size from environment variables and convert them to integers.
	maxMassagesToPoll, _ := strconv.Atoi(os.Getenv("MAX_MESSAGES"))
	maxWaitTimeToPoll, _ := strconv.Atoi(os.Getenv("MAX_WAIT_TIME"))
	sqsEndpoint := os.Getenv("SQS_ENDPOINT")
	encryptionKey := os.Getenv("ENCRYPTION_SECRET")

	// Initialize Logger
	logger, err := log.NewCustomLogger("logs")
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating logger with error %v", err.Error())}
		logger.Log(&lm)
	}

	lm := log.Message{Level: "INFO", Msg: "Logger initialized successfully"}
	logger.Log(&lm)

	// Initialize a new database connection.
	db := database.New(logger)
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating database failed with error %v", err.Error())}
		logger.Log(&lm)

		return
	}

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Database initilized sucessfully.")}
	logger.Log(&lm)

	// Initialize the ETL components.
	sqsClient, err := etl.NewSQSClient(logger, sqsEndpoint)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error initilizing sqs client: %v", err.Error())}
		logger.Log(&lm)

		return
	}

	extractor := etl.NewExtracter(logger, encryptionKey, sqsEndpoint, int32(maxMassagesToPoll), int32(maxWaitTimeToPoll))
	loader := etl.NewLoader(logger, dbConn, sqsClient)
	processor := etl.NewProcessor(logger, extractor, loader)

	// Start extraction and loading data
	processor.Worker()

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("SQS Client, Extractor, Loader, Processor initilized sucessfully.")}
	logger.Log(&lm)
}
