package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) AddAdmin(ctx context.Context, request *api.AddAdminRequest) (*api.AddAdminResponse, error) {
	if err := s.storage.AddAdmin(
		ctx,
		request.GetSecretId(),
		request.GetCallerAdminId(),
		request.GetCallerAdminPrivateKey(),
		request.GetNewAdminId(),
		request.GetNewAdminPublicKey(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.AddAdminResponse{}, nil
}
