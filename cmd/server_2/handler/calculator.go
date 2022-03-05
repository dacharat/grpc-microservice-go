package handler

import (
	"context"

	"github.com/dacharat/grpc-microservice-go/proto/calculator"
)

type Calculator struct{}

func (c *Calculator) Add(ctx context.Context, req *calculator.NumberReq) (*calculator.Result, error) {
	nums := req.GetNumbers()
	var result float32
	for _, n := range nums {
		result += float32(n)
	}
	return &calculator.Result{
		Output: result,
	}, nil
}
