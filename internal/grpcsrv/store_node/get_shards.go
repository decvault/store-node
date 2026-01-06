package store_node

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetShards(ctx context.Context, request *api.GetShardsRequest) (*api.GetShardsResponse, error) {
	secretID, err := uuid.Parse(request.SecretId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	shards, err := s.storage.GetShards(ctx, secretID, request.ShardIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetShardsResponse{
		Shards: shards,
	}, nil
}
