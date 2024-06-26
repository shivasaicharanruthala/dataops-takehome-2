package etl

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/shivasaicharanruthala/dataops-takehome-2/log"
	"github.com/shivasaicharanruthala/dataops-takehome-2/model"
)

type sqsActions struct {
	sqsClient   *sqs.Client
	logger      *log.CustomLogger
	sqsEndpoint string
}

func NewSQSClient(l *log.CustomLogger, sqsEndpoint string) (ISQSWrapper, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "ls-YAfoFAQe-5538-DENO-0592-KidEpibE1a33")))
	if err != nil {
		return nil, err
	}

	sc := sqs.NewFromConfig(sdkConfig, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(sqsEndpoint)
		o.Credentials = credentials.NewStaticCredentialsProvider("test", "test", "ls-YAfoFAQe-5538-DENO-0592-KidEpibE1a33")
	})
	return &sqsActions{
		sqsClient:   sc,
		logger:      l,
		sqsEndpoint: sqsEndpoint,
	}, nil
}

// DeleteMessages uses the DeleteMessageBatch action to delete a batch of messages from an Amazon SQS queue.
func (actor sqsActions) DeleteMessages(messages []*model.Response) error {
	entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for msgIndex, msg := range messages {
		entries[msgIndex].Id = aws.String(fmt.Sprintf("%v", msgIndex))
		entries[msgIndex].ReceiptHandle = &msg.ReceiptHandle
	}

	//TODO: Use context.Todo()
	_, err := actor.sqsClient.DeleteMessageBatch(context.TODO(), &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(actor.sqsEndpoint),
	})
	if err != nil {
		return err
	}

	return nil
}

// GetMessages uses the ReceiveMessage action to get messages from an Amazon SQS queue.
func (actor sqsActions) GetMessages(maxMessages int32, waitTime int32) ([]types.Message, error) {
	var messages []types.Message
	result, err := actor.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(actor.sqsEndpoint),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		fmt.Printf("Couldn't get messages from queue %v. Here's why: %v\n", actor.sqsEndpoint, err)
	} else {
		messages = result.Messages
	}
	return messages, err
}
