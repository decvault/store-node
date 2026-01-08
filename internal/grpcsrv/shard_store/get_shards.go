package shard_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/shard_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetShards(ctx context.Context, request *api.GetShardsRequest) (*api.GetShardsResponse, error) {
	shards, err := s.storage.GetShards(ctx, request.GetSecretId(), request.ShardIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetShardsResponse{
		Shards: shards,
	}, nil
}
