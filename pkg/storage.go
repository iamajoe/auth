package pkg

import (
	"github.com/iamajoe/auth/internal/storage"
)

type Connection = storage.Connection

func Dial(config *GlobalConfiguration) (*Connection, error) {
	return storage.Dial(config)
}
