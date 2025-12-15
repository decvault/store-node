package store_node

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
)

func (service *Service) PutShard(context.Context, *api.PutShardRequest) (*api.PutShardResponse, error) {
	return &api.PutShardResponse{}, nil
}
