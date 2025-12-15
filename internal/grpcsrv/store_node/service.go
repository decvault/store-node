package store_node

import (
	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
)

type Service struct {
	api.UnimplementedStoreNodeServer
}

func NewService() *Service {
	return &Service{}
}
