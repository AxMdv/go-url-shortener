package handlers

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"errors"
	"fmt"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/model"
	pb "github.com/AxMdv/go-url-shortener/internal/proto"
	"github.com/AxMdv/go-url-shortener/internal/storage"
	"github.com/AxMdv/go-url-shortener/pkg/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPCShortenerServer поддерживает все необходимые методы сервера.
type GRPCShortenerServer struct {
	// нужно встраивать тип pb.	pb.UnimplementedShortenerServer
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer

	shortenerService IShortenerService
	Config           config.Options
}

func NewGRPCShortenerServer(sS IShortenerService, Cfg config.Options) *GRPCShortenerServer {
	return &GRPCShortenerServer{shortenerService: sS, Config: Cfg}
}

// CreateShortURL creates short url from long url.
func (s *GRPCShortenerServer) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	var response pb.CreateShortURLResponse

	shortenedURL := s.shortenerService.ShortenLongURL([]byte(in.LongUrl))
	formedURL := &model.FormedURL{
		UUID:         auth.GetUUIDFromContext(ctx),
		ShortenedURL: shortenedURL,
		LongURL:      in.LongUrl,
	}
	err := s.shortenerService.CreateShortURL(formedURL)
	if err != nil {
		var duplicateErr *storage.AddURLError
		if errors.As(err, &duplicateErr) {
			response.ShortenedUrl = s.Config.ResponseResultAddr + duplicateErr.DuplicateValue
			return &response, status.Errorf(codes.AlreadyExists, "Url already exists %s", duplicateErr.DuplicateValue)
		}
		return nil, status.Errorf(codes.Internal, "Internal error %s", err.Error())
	}
	response.ShortenedUrl = s.Config.ResponseResultAddr + shortenedURL
	return &response, nil
}

// GetLongURL returns long url from shortened url
func (s *GRPCShortenerServer) GetLongURL(_ context.Context, in *pb.GetLongURLRequest) (*pb.GetLongURLResponse, error) {
	var response pb.GetLongURLResponse

	deleted, err := s.shortenerService.GetFlagByShortURL(in.ShortenedUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %s", err.Error())
	}
	if deleted {
		return nil, status.Errorf(codes.NotFound, "Already deleted url %s", in.ShortenedUrl)
	}

	longURL, err := s.shortenerService.GetLongURL(in.ShortenedUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %s", err.Error())
	}
	if longURL == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid url %s", in.ShortenedUrl)
	}
	response.LongUrl = longURL
	return &response, nil
}

// CheckDatabaseConnection checks if database is online.
func (s *GRPCShortenerServer) CheckDatabaseConnection(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {

	if s.Config.DataBaseDSN == "" {
		return nil, status.Errorf(codes.NotFound, "There is no DB")
	}
	err := s.shortenerService.PingDatabase()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %s", err.Error())
	}
	return nil, nil
}

// CreateShortURLBatch shortens batch of long urls.
func (s *GRPCShortenerServer) CreateShortURLBatch(ctx context.Context, in *pb.CreateShortURLBatchRequest) (*pb.CreateShortURLBatchResponse, error) {
	var response pb.CreateShortURLBatchResponse

	uuid := auth.GetUUIDFromContext(ctx)
	requestBatchList := make([]BatchOriginal, 0, len(in.OriginalBatches))

	for _, originalBatch := range in.OriginalBatches {
		batch := BatchOriginal{
			CorrelationID: originalBatch.CorrelationId,
			OriginalURL:   originalBatch.OriginalUrl,
		}
		requestBatchList = append(requestBatchList, batch)
	}
	requestBatch := RequestBatch{BatchList: requestBatchList}
	formedURL := requestBatch.ToFormed(uuid)
	err := s.shortenerService.CreateShortURLBatch(formedURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %s", err.Error())
	}
	for _, fu := range formedURL {
		shortenedBatch := &pb.CreateShortURLBatchResponse_BatchShortened{
			CorrelationId: fu.CorrelationID,
			ShortUrl:      fmt.Sprintf("%s/%s", s.Config.ResponseResultAddr, fu.ShortenedURL),
		}
		response.ShortenedBatches = append(response.ShortenedBatches, shortenedBatch)
	}
	return &response, nil
}

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func RegisterShortenerServer(s *grpc.Server, srvAPI *GRPCShortenerServer) {
	pb.RegisterShortenerServer(s, srvAPI)
}
