package PluginInterface

import (
	"github.com/EyciaZhou/msghub.go/interface"
	"time"
)

type PluginTask interface {
	FetchNew() ([]*Interface.Message, error)

	GetDelayBetweenCatchRound() time.Duration

	DumpTaskStatus() (Status []byte)

	Hash() string //as id, less equal 64 bytes

	Type() string
}
