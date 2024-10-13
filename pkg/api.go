package pkg

import (
	"net/http"

	"github.com/iamajoe/auth/internal/api"
)

type API = api.API
type Option = api.Option

func NewAPIWithVersionAndLogger(
	globalConfig *GlobalConfiguration,
	db *Connection,
	version string,
	loggerFn func(next http.Handler) http.Handler,
	opt ...Option,
) (*API, http.Handler) {
	return api.NewAPIWithVersionAndLogger(globalConfig, db, version, loggerFn, opt...)
}
