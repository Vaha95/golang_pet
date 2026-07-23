package grpcserver

import (
	"context"
	"errors"

	"github.com/Vaha95/golang_pet/api/shortenerpb"
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the ShortenerService gRPC service.
type Server struct {
	shortenerpb.UnimplementedShortenerServiceServer
	cfg     config.StorageConfig
	auditCh chan DTO.BaseAuditItem
}

func (s *Server) sendAudit(item DTO.BaseAuditItem) {
	select {
	case s.auditCh <- item:
	default:
	}
}

// NewServer creates a new gRPC server.
func NewServer(cfg config.StorageConfig, auditCh chan DTO.BaseAuditItem) *Server {
	return &Server{cfg: cfg, auditCh: auditCh}
}

// toGRPCStatus maps repository-level errors to gRPC status codes.
func toGRPCStatus(err error) error {
	switch {
	case errors.Is(err, repository.ErrorShortURLKeyNotFound):
		return status.Error(codes.NotFound, "URL is not found")
	case errors.Is(err, strategy.ErrorUrlNotFound):
		return status.Error(codes.NotFound, "URL is not found")
	case errors.Is(err, repository.ErrorURLByUserNotFound):
		return status.Error(codes.NotFound, "URL is not found")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// toGRPCStatusSave maps saveurl errors to gRPC status codes.
func toGRPCStatusSave(err error) error {
	switch {
	case errors.Is(err, saveurl.ErrorParseRequestURI):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, saveurl.ErrorSaveToStorage):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.InvalidArgument, err.Error())
	}
}

func (s *Server) Shorten(ctx context.Context, req *shortenerpb.ShortenRequest) (*shortenerpb.ShortenResponse, error) {
	path, err := saveurl.SaveURL(s.cfg, req.Url)
	if errors.Is(err, saveurl.ErrorUrlAlreadyExists) {
		s.sendAudit(DTO.CreateBaseAuditItemFollow(req.Url, s.cfg.GetUserId()))
		return &shortenerpb.ShortenResponse{ShortUrl: path}, nil
	}
	if err != nil {
		return nil, toGRPCStatusSave(err)
	}

	s.sendAudit(DTO.CreateBaseAuditItemFollow(req.Url, s.cfg.GetUserId()))

	return &shortenerpb.ShortenResponse{ShortUrl: path}, nil
}

func (s *Server) Resolve(ctx context.Context, req *shortenerpb.ResolveRequest) (*shortenerpb.ResolveResponse, error) {
	urlStorage := strategy.GetStrategy(s.cfg)
	data, err := urlStorage.Get(req.Id)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	if data == nil || data.URL == "" {
		return nil, status.Error(codes.Internal, "URL is empty")
	}

	s.sendAudit(DTO.CreateBaseAuditItemFollow(data.URL, s.cfg.GetUserId()))

	return &shortenerpb.ResolveResponse{
		Url:     data.URL,
		Deleted: data.DeletedAt != nil,
	}, nil
}

func (s *Server) GetUserURLs(ctx context.Context, req *shortenerpb.GetUserURLsRequest) (*shortenerpb.GetUserURLsResponse, error) {
	urlStorage := strategy.GetStrategy(s.cfg)
	data, err := urlStorage.GetByUser(s.cfg.GetUserId())
	if errors.Is(err, repository.ErrorURLByUserNotFound) {
		return &shortenerpb.GetUserURLsResponse{}, nil
	}
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	if len(data) == 0 {
		return &shortenerpb.GetUserURLsResponse{}, nil
	}

	var urls []*shortenerpb.UserURL
	for _, v := range data {
		urls = append(urls, &shortenerpb.UserURL{
			Short:   v.Short,
			Url:     v.URL,
			Deleted: v.DeletedAt != nil,
		})
	}

	return &shortenerpb.GetUserURLsResponse{Urls: urls}, nil
}
