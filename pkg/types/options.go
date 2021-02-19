package types

import (
	"github.com/MaxPolarfox/goTools/mongoDB"
)

type Options struct {
	Port int `json:"port"`
	ServiceName string `json:"serviceName"`
	DB Collections `json:"db"`
}

type Collections struct {
	Tasks mongoDB.Options `json:"tasks"`
}
