package test

import (
	"github.com/iamajoe/auth/internal/conf"
	"github.com/iamajoe/auth/internal/storage"
)

func SetupDBConnection(globalConfig *conf.GlobalConfiguration) (*storage.Connection, error) {
	return storage.Dial(globalConfig)
}
