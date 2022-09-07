// Code generated by https://github.com/meerkat-io/disorder; DO NOT EDIT.
package sub

import (
	"fmt"
	"github.com/meerkat-io/disorder"
	"github.com/meerkat-io/disorder/rpc"
	"github.com/meerkat-io/disorder/rpc/code"
)

type mathServiceHandler func(*rpc.Context, *disorder.Decoder) (interface{}, *rpc.Error)

type MathService interface {
	Increase(*rpc.Context, int32) (int32, *rpc.Error)
}

func NewMathServiceClient(client *rpc.Client) MathService {
	return &mathServiceClient{
		name:   "math_service",
		client: client,
	}
}

type mathServiceClient struct {
	name   string
	client *rpc.Client
}

func (c *mathServiceClient) Increase(context *rpc.Context, request int32) (int32, *rpc.Error) {
	var response int32
	err := c.client.Send(context, c.name, "increase", request, &response)
	return response, err
}

type mathServiceServer struct {
	name    string
	service MathService
	methods map[string]mathServiceHandler
}

func RegisterMathServiceServer(s *rpc.Server, service MathService) {
	server := &mathServiceServer{
		name:    "math_service",
		service: service,
	}
	server.methods = map[string]mathServiceHandler{
		"increase": server.increase,
	}
	s.RegisterService("math_service", server)
}

func (s *mathServiceServer) Handle(context *rpc.Context, method string, d *disorder.Decoder) (interface{}, *rpc.Error) {
	handler, ok := s.methods[method]
	if ok {
		return handler(context, d)
	}
	return nil, &rpc.Error{
		Code:  code.Unimplemented,
		Error: fmt.Errorf("unimplemented method \"%s\" under service \"%s\"", method, s.name),
	}
}

func (s *mathServiceServer) increase(context *rpc.Context, d *disorder.Decoder) (interface{}, *rpc.Error) {
	var request int32
	err := d.Decode(&request)
	if err != nil {
		return nil, &rpc.Error{
			Code:  code.InvalidRequest,
			Error: err,
		}
	}
	response, rpcErr := s.service.Increase(context, request)
	if rpcErr != nil {
		return nil, rpcErr
	}
	return response, nil
}
