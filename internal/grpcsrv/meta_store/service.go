package meta_store

import (
	api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	"github.com/decvault/store-node/internal/pkg/storage"
)

type Service struct {
	api.UnimplementedMetaStoreServer

	storage storage.MetaStorage
}

func NewService(
	storage storage.MetaStorage,
) *Service {
	return &Service{
		storage: storage,
	}
}
