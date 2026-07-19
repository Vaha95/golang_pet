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

// NewServer creates a new gRPC server.
func NewServer(cfg config.StorageConfig, auditCh chan DTO.BaseAuditItem) *Server {
	return &Server{cfg: cfg, auditCh: auditCh}
}

func (s *Server) Shorten(ctx context.Context, req *shortenerpb.ShortenRequest) (*shortenerpb.ShortenResponse, error) {
	path, err := saveurl.SaveURL(s.cfg, req.Url)
	if err != nil {
		if errors.Is(err, saveurl.ErrorParseRequestURI) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, saveurl.ErrorUrlAlreadyExists) {
			s.auditCh <- DTO.CreateBaseAuditItemFollow(req.Url, s.cfg.GetUserId())
			return &shortenerpb.ShortenResponse{ShortUrl: path}, nil
		}

		if errors.Is(err, saveurl.ErrorSaveToStorage) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s.auditCh <- DTO.CreateBaseAuditItemFollow(req.Url, s.cfg.GetUserId())

	return &shortenerpb.ShortenResponse{ShortUrl: path}, nil
}

func (s *Server) Resolve(ctx context.Context, req *shortenerpb.ResolveRequest) (*shortenerpb.ResolveResponse, error) {
	s2 := strategy.GetStrategy(s.cfg)
	data, err := s2.Get(req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrorShortURLKeyNotFound) || errors.Is(err, strategy.ErrorUrlNotFound) {
			return nil, status.Error(codes.NotFound, "URL is not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	if data == nil || data.URL == "" {
		return nil, status.Error(codes.Internal, "URL is empty")
	}

	s.auditCh <- DTO.CreateBaseAuditItemFollow(data.URL, s.cfg.GetUserId())

	return &shortenerpb.ResolveResponse{
		Url:     data.URL,
		Deleted: data.DeletedAt != nil,
	}, nil
}

func (s *Server) GetUserURLs(ctx context.Context, req *shortenerpb.GetUserURLsRequest) (*shortenerpb.GetUserURLsResponse, error) {
	s2 := strategy.GetStrategy(s.cfg)
	data, err := s2.GetByUser(s.cfg.GetUserId())
	if err != nil {
		if errors.Is(err, repository.ErrorURLByUserNotFound) {
			return &shortenerpb.GetUserURLsResponse{}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
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
