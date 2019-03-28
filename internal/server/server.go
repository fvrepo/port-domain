package server

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"

	"github.com/port-domain/internal/models"
	portApi "github.com/port-domain/pkg/grpcapi/port"
)

var l = logrus.New()

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
	p := grpcToPortModel(request)
	if err := s.controller.SavePort(ctx, p); err != nil {
		l.WithError(err).Error()
		return nil, status.Error(codes.Internal, errors.WithStack(err).Error())
	}
	return &portApi.SavePortResponse{}, nil
}

func (s *Server) GetAllPorts(ctx context.Context, request *portApi.GetAllPortsRequest) (*portApi.GetAllPortsResponse, error) {
	pm, err := s.controller.GetAllPorts(ctx, int(request.GetLimit()), int(request.GetLimit()))
	if err != nil {
		l.WithError(err).Error()
		return nil, status.Error(codes.Internal, errors.WithStack(err).Error())
	}

	portMap := &portApi.PortMap{Port: modelToGrpc(pm)}

	return &portApi.GetAllPortsResponse{Ports: portMap}, nil
}

func grpcToPortModel(req *portApi.SavePortRequest) *models.Port {
	return &models.Port{
		ID:          req.Id,
		Country:     req.Details.Country,
		City:        req.Details.City,
		Code:        req.Details.Code,
		Alias:       req.Details.Alias,
		Regions:     req.Details.Regions,
		Coordinates: req.Details.Coordinates,
		Province:    req.Details.Province,
		Timezone:    req.Details.Timezone,
		Unlocs:      req.Details.Unlocs,
	}
}

func modelToGrpc(ports []*models.Port) map[string]*portApi.PortDetails {
	pd := make(map[string]*portApi.PortDetails)
	for _, p := range ports {
		pd[p.ID] = &portApi.PortDetails{
			Country:     p.Country,
			City:        p.City,
			Code:        p.Code,
			Alias:       p.Alias,
			Regions:     p.Regions,
			Coordinates: p.Coordinates,
			Province:    p.Province,
			Timezone:    p.Timezone,
			Unlocs:      p.Unlocs,
		}
	}
	return pd
}
