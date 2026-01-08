package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveReader(ctx context.Context, request *api.RemoveReaderRequest) (*api.RemoveReaderResponse, error) {
	if err := s.storage.RemoveReader(
		ctx,
		request.GetSecretId(),
		request.GetCallerAdminId(),
		request.GetCallerAdminPrivateKey(),
		request.GetRemovedReaderId(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.RemoveReaderResponse{}, nil
}
