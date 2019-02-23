package endpoint

import (
	"context"

	"github.com/eirsyl/flexit/examples/simple/pkg/service"

	flexitendpoint "github.com/eirsyl/flexit/endpoint"
	"github.com/eirsyl/flexit/examples/simple/pb"
)

// Set defines the service set
type Set struct {
	AddEndpoint      flexitendpoint.Endpoint
	SubtractEndpoint flexitendpoint.Endpoint
}

var _ service.Service = &Set{}

// Add implements the add function of the service definition
func (s *Set) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	resp, err := s.AddEndpoint(ctx, r)
	response := resp.(*pb.AddResponse)
	return response, err
}

// Subtract implements the subtract function of the service definition
func (s *Set) Subtract(ctx context.Context, r *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	resp, err := s.SubtractEndpoint(ctx, r)
	response := resp.(*pb.SubtractResponse)
	return response, err
}
