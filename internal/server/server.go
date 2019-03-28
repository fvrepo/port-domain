package server

import (
	"context"

	"github.com/port-domain/internal/models"

	portApi "github.com/port-domain/pkg/grpcapi/port"
)

type Controller interface {
	SavePort(ctx context.Context, p *models.Port) error
	GetAllPorts(ctx context.Context, limit, skip int) ([]*models.Port, error)
}

type Server struct {
	controller Controller
}

func New(con Controller) *Server {
	return &Server{controller: con}
}

func (s *Server) SavePort(ctx context.Context, request *portApi.SavePortRequest) (*portApi.SavePortResponse, error) {

	return &portApi.SavePortResponse{}, nil
}

func (s *Server) GetAllPorts(ctx context.Context, request *portApi.GetAllPortsRequest) (*portApi.GetAllPortsResponse, error) {
	return &portApi.GetAllPortsResponse{}, nil
}
