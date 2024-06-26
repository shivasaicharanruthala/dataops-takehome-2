package etl

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/shivasaicharanruthala/dataops-takehome-2/model"
)

type IProcessor interface {
	Worker()
}

type IExtractor interface {
	FetchDataFromSQS() ([]model.Response, error)
}

type ILoader interface {
	BatchInsert(responses []model.Response) error
	SequentialInsert(response []model.Response) error
}

type ISQSWrapper interface {
	GetMessages(maxMessages int32, waitTime int32) ([]types.Message, error)
	DeleteMessages(messages []*model.Response) error
}
