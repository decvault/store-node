package shard_store

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/shard_store/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) PutShards(ctx context.Context, request *api.PutShardsRequest) (*api.PutShardsResponse, error) {
	if err := s.storage.SaveShards(ctx, request.GetSecretId(), request.Shards); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.PutShardsResponse{}, nil
}
