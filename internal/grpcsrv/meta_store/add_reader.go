package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) AddReader(ctx context.Context, request *api.AddReaderRequest) (*api.AddReaderResponse, error) {
	if err := s.storage.AddReader(
		ctx,
		request.GetSecretId(),
		request.GetCallerAdminId(),
		request.GetCallerAdminPrivateKey(),
		request.GetNewReaderId(),
		request.GetNewReaderPublicKey(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.AddReaderResponse{}, nil
}
