package shard_store

import (
	api "github.com/decvault/store-node/internal/pb/github.com/decvault/shard_store/api"
	"github.com/decvault/store-node/internal/pkg/storage"
)

type Service struct {
	api.UnimplementedShardStoreServer

	storage storage.ShardStorage
}

func NewService(
	storage storage.ShardStorage,
) *Service {
	return &Service{
		storage: storage,
	}
}
