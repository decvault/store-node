package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveAdmin(ctx context.Context, request *api.RemoveAdminRequest) (*api.RemoveAdminResponse, error) {
	if err := s.storage.RemoveAdmin(
		ctx,
		request.GetSecretId(),
		request.GetCallerAdminId(),
		request.GetCallerAdminPrivateKey(),
		request.GetRemovedAdminId(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.RemoveAdminResponse{}, nil
}
