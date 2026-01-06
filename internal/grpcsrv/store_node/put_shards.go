package store_node

import (
	"context"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Service) PutShards(ctx context.Context, request *api.PutShardsRequest) (*api.PutShardsResponse, error) {
	secretID, err := uuid.Parse(request.SecretId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = s.storage.SaveShards(ctx, secretID, request.Shards)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &api.PutShardsResponse{}, nil
}
