package client

import (
	"context"
	"sync"

	"github.com/PaulYakow/metrics-track/cmd/agent/config"
)

type ISender interface {
	Run(context.Context, *sync.WaitGroup, *config.Config)
}
