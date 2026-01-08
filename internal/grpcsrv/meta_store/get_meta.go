package meta_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetMeta(ctx context.Context, request *api.GetMetaRequest) (*api.GetMetaResponse, error) {
	dek, threshold, shards, err := s.storage.GetMeta(ctx, request.GetSecretId(), request.GetCallerId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetMetaResponse{
		EncryptedDek: dek,
		Shards:       shards,
		Threshold:    threshold,
	}, nil
}
