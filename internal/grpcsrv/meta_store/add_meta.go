package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) AddMeta(ctx context.Context, request *api.AddMetaRequest) (*api.AddMetaResponse, error) {
	if err := s.storage.AddMeta(
		ctx,
		request.GetSecretId(),
		request.GetAdminId(),
		request.GetAdminPublicKey(),
		request.GetDek(),
		request.GetThreshold(),
		request.GetShards(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.AddMetaResponse{}, nil
}
