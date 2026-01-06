package store_node

import (
	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
	"github.com/decvault/store-node/internal/pkg/storage"
)

type Service struct {
	api.UnimplementedStoreNodeServer

	storage storage.ShardStorage
}

func NewService(
	storage storage.ShardStorage,
) *Service {
	return &Service{
		storage: storage,
	}
}
