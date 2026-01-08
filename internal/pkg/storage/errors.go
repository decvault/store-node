package storage

import (
	"errors"
)

var (
	ErrMetaNotFound                        = errors.New("meta not found")
	ErrMetaAlreadyExists                   = errors.New("meta already exists")
	ErrNilMeta                             = errors.New("meta can not be nil")
	ErrCallerHasNoReadAccess               = errors.New("caller has no read access")
	ErrCallerIsNotInAdminList              = errors.New("caller is not in admin list")
	ErrAdminAlreadyExists                  = errors.New("admin already exists")
	ErrAdminDoesNotExist                   = errors.New("admin does not exist")
	ErrReaderAlreadyExists                 = errors.New("reader already exists")
	ErrReaderDoesNotExist                  = errors.New("reader does not exist")
	ErrInvalidShardsThresholdConfiguration = errors.New("invalid shards/threshold configuration")
)
