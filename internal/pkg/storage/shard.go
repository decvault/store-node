package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

type ShardStorage interface {
	GetShards(ctx context.Context, secretID string, shardIDs []uint32) (map[uint32][]byte, error)
	SaveShards(ctx context.Context, secretID string, shards map[uint32][]byte) error
}

type shardStorage struct {
	db *badger.DB
}

func NewShardStorage(db *badger.DB) ShardStorage {
	return &shardStorage{
		db: db,
	}
}

func (s *shardStorage) GetShards(_ context.Context, secretID string, shardIDs []uint32) (map[uint32][]byte, error) {
	shards := make(map[uint32][]byte)
	err := s.db.View(func(txn *badger.Txn) error {
		for _, shardID := range shardIDs {
			key := shardKey(secretID, shardID)
			item, err := txn.Get(key)
			if errors.Is(err, badger.ErrKeyNotFound) {
				continue
			} else if err != nil {
				return errors.WithStack(err)
			}

			raw, err := item.ValueCopy(nil)
			if err != nil {
				return errors.WithStack(err)
			}

			var value shardValue
			if err := json.Unmarshal(raw, &value); err != nil {
				return errors.WithStack(err)
			}

			shards[shardID] = value.Data
		}

		return nil
	})

	return shards, errors.WithStack(err)
}

func (s *shardStorage) SaveShards(_ context.Context, secretID string, shards map[uint32][]byte) error {
	currentTimestamp := time.Now().UnixMilli()
	return s.db.Update(func(txn *badger.Txn) error {
		for shardID, data := range shards {
			value, err := json.Marshal(shardValue{
				Data:               data,
				CreatedAtTimestamp: currentTimestamp,
			})

			if err != nil {
				return errors.WithStack(err)
			}

			if err := txn.Set(shardKey(secretID, shardID), value); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

func shardKey(secretID string, shardID uint32) []byte {
	return []byte(fmt.Sprintf("secret:%s:%d", secretID, shardID))
}

type shardValue struct {
	Data               []byte `json:"data"`
	CreatedAtTimestamp int64  `json:"created_at"`
}
