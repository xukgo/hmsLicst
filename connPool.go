package hmsLicenter

import (
	"fmt"
	"github.com/flyaways/pool"
	"google.golang.org/grpc"
	"time"
)

func initPool(addr string) (*pool.GRPCPool, error) {
	options := &pool.Options{
		InitTargets:  []string{addr},
		InitCap:      2,
		MaxCap:       5,
		DialTimeout:  time.Second * 3,
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 2,
	}

	p, err := pool.NewGRPCPool(options, grpc.WithInsecure()) //for grpc

	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("grpc conn pool is nil")
	}

	return p, nil
}
