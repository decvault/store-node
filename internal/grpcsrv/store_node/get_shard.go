package store_node

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
)

func (service *Service) GetShard(context.Context, *api.GetShardRequest) (*api.GetShardResponse, error) {
	return &api.GetShardResponse{}, nil
}
