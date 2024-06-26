package etl

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome-2/log"
	"github.com/shivasaicharanruthala/dataops-takehome-2/model"
	"io"
	"net/http"
)

type extract struct {
	httpClient    *http.Client
	logger        *log.CustomLogger
	sqsEndpoint   string
	noOfMessages  int32
	waitTimeInSec int32
	encryptionKey string
}

// NewExtracter creates a new instance of the Extractor and initializes it with the SQS endpoint from environment variables.
func NewExtracter(logger *log.CustomLogger, encryptionKey string, sqsEndpoint string, noOfMessages, waitTimeInSec int32) IExtractor {
	return &extract{
		httpClient:    new(http.Client),
		logger:        logger,
		sqsEndpoint:   sqsEndpoint,
		noOfMessages:  noOfMessages,
		waitTimeInSec: waitTimeInSec,
		encryptionKey: encryptionKey,
	}
}

// FetchDataFromSQS fetches data from the SQS endpoint, processes the response, and returns a model.Response.
func (ex extract) FetchDataFromSQS() ([]model.Response, error) {
	endpoint := ex.sqsEndpoint + fmt.Sprintf("&MaxNumberOfMessages=%v&WaitTimeSeconds=%v", ex.noOfMessages, ex.waitTimeInSec)

	// Create a new GET request to the SQS endpoint.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error creating request to sqs enpoint: %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Send the request and receive the response.
	resp, err := ex.httpClient.Do(req)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error sending request to sqs enpoint: %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error reading response from sqs enpoint : %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Unmarshal the XML response into the ReceiveMessageResponse struct.
	var sqsMessageResponse model.ReceiveMessageResponse
	err = xml.Unmarshal(body, &sqsMessageResponse)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error unmarshalling XML response from sqs enpoint : %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Initialize a new Response struct.
	var msglist []model.Response
	if sqsMessageResponse.ReceiveMessageResult.Message != nil && len(sqsMessageResponse.ReceiveMessageResult.Message) > 0 {
		for _, msg := range sqsMessageResponse.ReceiveMessageResult.Message {
			var res model.Response

			// Unmarshal the JSON body of the SQS message into the Response struct.
			err = json.Unmarshal([]byte(msg.Body), &res)
			if err != nil {
				lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error unmarshalling JSON body from XML response from sqs enpoint : %v", err.Error())}
				ex.logger.Log(&lm)

				fmt.Println("Error unmarshalling JSON body:", err)
				return nil, err
			}

			if msg.MessageId != nil && res.UserID != nil {
				// Set additional data from the SQS message response into the Response struct.
				res.SetData(sqsMessageResponse.ResponseMetadata.RequestId, &msg)

				// Mask sensitive data in the Response struct.
				err = res.MaskBody(ex.encryptionKey)
				if err != nil {
					return nil, err
				}

				msglist = append(msglist, res)
			}
		}
	}

	// Return the populated Response struct.
	return msglist, nil
}
