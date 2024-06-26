package store

import "github.com/shivasaicharanruthala/dataops-takehome-2/model"

type Login interface {
	Get(filter *model.Filter) ([]model.Response, error)
}
