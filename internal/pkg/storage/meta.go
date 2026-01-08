package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

type secretMeta struct {
	Admins       map[string]struct{} `json:"admins"`
	Readers      map[string]struct{} `json:"readers"`
	EncryptedDEK map[string][]byte   `json:"encrypted_dek"`
	Threshold    int64               `json:"threshold"`
	Shards       int64               `json:"shards"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type MetaStorage interface {
	GetMeta(ctx context.Context, secretID string, callerID string) (dek []byte, threshold int64, shards int64, err error)
	AddMeta(ctx context.Context, secretID string, adminID string, adminPublicKey []byte, dek []byte, threshold int64, shards int64) error
	AddAdmin(ctx context.Context, secretID string, callerAdminID string, callerAdminPrivateKey []byte, newAdminID string, newAdminPublicKey []byte) error
	RemoveAdmin(ctx context.Context, secretID string, callerAdminID string, callerAdminPrivateKey []byte, removedAdminID string) error
	AddReader(ctx context.Context, secretID string, callerAdminID string, callerAdminPrivateKey []byte, newReaderID string, newReaderPublicKey []byte) error
	RemoveReader(ctx context.Context, secretID string, callerAdminID string, callerAdminPrivateKey []byte, removedReaderID string) error
}

type metaStorage struct {
	db *badger.DB
}

func NewMetaStorage(db *badger.DB) MetaStorage {
	return &metaStorage{
		db: db,
	}
}

func (m *metaStorage) GetMeta(
	ctx context.Context,
	secretID string,
	callerID string,
) ([]byte, int64, int64, error) {
	meta, err := m.getMeta(ctx, secretID)
	if err != nil {
		return nil, 0, 0, err
	}

	_, adminOK := meta.Admins[callerID]
	_, readerOK := meta.Readers[callerID]
	if !adminOK || !readerOK {
		return nil, 0, 0, errors.WithStack(ErrCallerHasNoReadAccess)
	}

	return meta.EncryptedDEK[callerID], meta.Threshold, meta.Shards, nil
}

func (m *metaStorage) AddMeta(
	_ context.Context,
	secretID string,
	adminID string,
	adminPublicKey []byte,
	dek []byte,
	threshold int64,
	shards int64,
) error {
	if threshold <= 0 || shards <= 0 || threshold > shards {
		return errors.WithStack(ErrInvalidShardsThresholdConfiguration)
	}

	encryptedDEK, err := encryptDEK(dek, adminPublicKey)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	meta := &secretMeta{
		Admins: map[string]struct{}{
			adminID: {},
		},
		Readers: make(map[string]struct{}),
		EncryptedDEK: map[string][]byte{
			adminID: encryptedDEK,
		},
		Threshold: threshold,
		Shards:    shards,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	return m.db.Update(func(txn *badger.Txn) error {
		key := metaKey(secretID)
		_, err := txn.Get(key)
		if err == nil {
			return errors.WithStack(ErrMetaAlreadyExists)
		} else if !errors.Is(err, badger.ErrKeyNotFound) {
			return errors.WithStack(err)
		}

		data, err := json.Marshal(meta)
		if err != nil {
			return errors.WithStack(err)
		}

		return errors.WithStack(txn.Set(key, data))
	})
}

func (m *metaStorage) AddAdmin(
	ctx context.Context,
	secretID string,
	callerAdminID string,
	callerAdminPrivateKey []byte,
	newAdminID string,
	newAdminPublicKey []byte,
) error {
	meta, err := m.getMeta(ctx, secretID)
	if err != nil {
		return err
	}

	if _, ok := meta.Admins[callerAdminID]; !ok {
		return errors.WithStack(ErrCallerIsNotInAdminList)
	}

	if _, ok := meta.Admins[newAdminID]; ok {
		return errors.WithStack(ErrAdminAlreadyExists)
	}

	dek, err := decryptDEK(meta.EncryptedDEK[callerAdminID], callerAdminPrivateKey)
	if err != nil {
		return err
	}

	newEncryptedDEK, err := encryptDEK(dek, newAdminPublicKey)
	if err != nil {
		return err
	}

	meta.Admins[newAdminID] = struct{}{}
	meta.EncryptedDEK[newAdminID] = newEncryptedDEK

	return m.updateMeta(ctx, secretID, meta)
}

func (m *metaStorage) RemoveAdmin(
	ctx context.Context,
	secretID string,
	callerAdminID string,
	callerAdminPrivateKey []byte,
	removedAdminID string,
) error {
	meta, err := m.getMeta(ctx, secretID)
	if err != nil {
		return err
	}

	if _, ok := meta.Admins[callerAdminID]; !ok {
		return errors.WithStack(ErrCallerIsNotInAdminList)
	}

	if _, ok := meta.Admins[removedAdminID]; !ok {
		return errors.WithStack(ErrAdminDoesNotExist)
	}

	_, err = decryptDEK(meta.EncryptedDEK[callerAdminID], callerAdminPrivateKey)
	if err != nil {
		return err
	}

	delete(meta.Admins, removedAdminID)
	delete(meta.EncryptedDEK, removedAdminID)

	return m.updateMeta(ctx, secretID, meta)
}

func (m *metaStorage) AddReader(
	ctx context.Context,
	secretID string,
	callerAdminID string,
	callerAdminPrivateKey []byte,
	newReaderID string,
	newReaderPublicKey []byte,
) error {
	meta, err := m.getMeta(ctx, secretID)
	if err != nil {
		return err
	}

	if _, ok := meta.Admins[callerAdminID]; !ok {
		return errors.WithStack(ErrCallerIsNotInAdminList)
	}

	if _, ok := meta.Readers[newReaderID]; ok {
		return errors.WithStack(ErrReaderAlreadyExists)
	}

	dek, err := decryptDEK(meta.EncryptedDEK[callerAdminID], callerAdminPrivateKey)
	if err != nil {
		return err
	}

	newEncryptedDEK, err := encryptDEK(dek, newReaderPublicKey)
	if err != nil {
		return err
	}

	meta.Readers[newReaderID] = struct{}{}
	meta.EncryptedDEK[newReaderID] = newEncryptedDEK

	return m.updateMeta(ctx, secretID, meta)
}

func (m *metaStorage) RemoveReader(
	ctx context.Context,
	secretID string,
	callerAdminID string,
	callerAdminPrivateKey []byte,
	removedReaderID string,
) error {
	meta, err := m.getMeta(ctx, secretID)
	if err != nil {
		return err
	}

	if _, ok := meta.Admins[callerAdminID]; !ok {
		return errors.WithStack(ErrCallerIsNotInAdminList)
	}

	if _, ok := meta.Admins[removedReaderID]; !ok {
		return errors.WithStack(ErrReaderDoesNotExist)
	}

	_, err = decryptDEK(meta.EncryptedDEK[callerAdminID], callerAdminPrivateKey)
	if err != nil {
		return err
	}

	delete(meta.Readers, removedReaderID)
	delete(meta.EncryptedDEK, removedReaderID)

	return m.updateMeta(ctx, secretID, meta)
}

func (m *metaStorage) getMeta(
	_ context.Context,
	secretID string,
) (*secretMeta, error) {
	var meta secretMeta
	err := m.db.View(func(txn *badger.Txn) error {
		key := metaKey(secretID)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(value []byte) error {
			return json.Unmarshal(value, &meta)
		})
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, errors.WithStack(ErrMetaNotFound)
		}

		return nil, errors.WithStack(err)
	}

	return &meta, nil
}

func (m *metaStorage) updateMeta(
	_ context.Context,
	secretID string,
	meta *secretMeta,
) error {
	if meta == nil {
		return ErrNilMeta
	}

	return m.db.Update(func(txn *badger.Txn) error {
		key := metaKey(secretID)
		_, err := txn.Get(key)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return errors.WithStack(ErrMetaNotFound)
			}

			return errors.WithStack(err)
		}

		meta.UpdatedAt = time.Now()
		data, err := json.Marshal(meta)
		if err != nil {
			return err
		}

		return errors.WithStack(txn.Set(key, data))
	})
}

func encryptDEK(dek []byte, rawPublicKey []byte) ([]byte, error) {
	publicKey, err := x509.ParsePKIXPublicKey(rawPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}

	publicRSA, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicRSA, dek, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt DEK")
	}

	return ciphertext, nil
}

func decryptDEK(encryptedDEK []byte, rawPrivateKey []byte) ([]byte, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(rawPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedDEK, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt DEK")
	}

	return plaintext, nil
}

func metaKey(secretID string) []byte {
	return []byte(fmt.Sprintf("meta:%s", secretID))
}
