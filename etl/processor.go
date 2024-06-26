package etl

import (
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome-2/log"
	"os"
	"strconv"
	"time"
)

type transformer struct {
	logger    *log.CustomLogger
	extractor IExtractor
	loader    ILoader
}

// NewProcessor creates a new instance of the Processor with the given extractor and loader.
func NewProcessor(logger *log.CustomLogger, extractor IExtractor, loader ILoader) IProcessor {
	return &transformer{
		logger:    logger,
		extractor: extractor,
		loader:    loader,
	}
}

// Worker is a function that continuously calls the API to fetch data and sends the result to a channel.
func (p *transformer) Worker() {
	// Retrieve maximum allowed empty responses and maximum consecutive empty responses from environment variables
	maxEmptyResponses, _ := strconv.Atoi(os.Getenv("MAX_NO_RESPONSES"))
	maxConsecutiveEmptyResponses, _ := strconv.Atoi(os.Getenv("MAX_CONSECUTIVE_NO_RESPONSES"))
	initialWaitTime := time.Second
	emptyResponseCount := 0
	waitTime := initialWaitTime

	for {
		response, err := p.extractor.FetchDataFromSQS()
		if err != nil {
			lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error fetching data: %v", err.Error())}
			p.logger.Log(&lm)

			continue
		}

		if len(response) > 0 {
			err = p.loader.BatchInsert(response)
			if err != nil {
				lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error inserting batch: %v", err.Error())}
				p.logger.Log(&lm)

			}

			emptyResponseCount = 0     // Reset the empty response counter
			waitTime = initialWaitTime // Reset the wait time
		} else {
			lm := log.Message{Level: "INFO", Msg: fmt.Sprintf("Received empty response")}
			p.logger.Log(&lm)

			emptyResponseCount++
			if emptyResponseCount >= maxEmptyResponses {
				lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Waiting for %v due to consecutive empty responses", waitTime)}
				p.logger.Log(&lm)

				time.Sleep(waitTime)        // Wait for the specified time before retrying
				waitTime += initialWaitTime // Increase the wait time linearly
			}

			if emptyResponseCount >= maxConsecutiveEmptyResponses {
				lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Reached max consecutive empty responses, canceling context")}
				p.logger.Log(&lm)

				return
			}
		}
	}
}
