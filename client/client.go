package client

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/iamajoe/auth/pkg"
	"github.com/iamajoe/auth/pkg/psqlutils"
)

// REF: this code is mostly copied from supabase/auth

type authLogger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type AuthAppConfig struct {
	DbHost               string
	DbPort               int
	DbUsername           string
	DbPassword           string
	DbName               string
	DbSslmode            bool
	DbAuthNamespace      string
	DbMigrationNamespace string
}

type Auth struct {
	db        *pkg.Connection
	handler   http.Handler
	log       authLogger
	pkgConfig *pkg.GlobalConfiguration
	appConfig AuthAppConfig
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewClient(
	ctx context.Context,
	pkgConf *pkg.GlobalConfiguration,
	log authLogger,
	httpLogger func(next http.Handler) http.Handler,
	version string,

) (*Auth, error) {
	log.Info("connecting auth to database")

	db, err := pkg.Dial(pkgConf)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %+v", err)
	}
	defer db.Close()

	opts := []pkg.Option{
		pkg.NewLimiterOptions(pkgConf),
	}
	_, handler := pkg.NewAPIWithVersionAndLogger(
		pkgConf,
		db,
		version,
		httpLogger,
		// logger.NewHTTPHandler(log),
		opts...,
	)

	return &Auth{
		db:        db,
		handler:   handler,
		pkgConfig: pkgConf,
		log:       log,
	}, nil

}

func NewClientWithConfig(
	ctx context.Context,
	conf AuthAppConfig,
	pkgConf *pkg.GlobalConfiguration,
	log authLogger,
	httpLogger func(next http.Handler) http.Handler,
	version string,
) (*Auth, error) {
	client, err := NewClient(ctx, pkgConf, log, httpLogger, version)
	if err != nil {
		return nil, err
	}

	client.appConfig = conf

	return client, nil
}

func (g *Auth) Close() error {
	if g.db == nil || g.db.Store == nil {
		return nil
	}

	return g.db.Close()
}

func (g *Auth) MigrateUp() error {
	g.log.Info("processing auth migrations...")

	err := psqlutils.MigrateUp(psqlutils.ConnectionDetails{
		// MigAppName: "gi_auth_migrations",
		MigAppName: g.appConfig.DbMigrationNamespace,
		Host:       g.appConfig.DbHost,
		Port:       g.appConfig.DbPort,
		Username:   g.appConfig.DbUsername,
		Password:   g.appConfig.DbPassword,
		Dbname:     g.appConfig.DbName,
		Sslmode:    g.appConfig.DbSslmode,
		Options: map[string]string{
			"migration_table_name": "schema_" + g.appConfig.DbMigrationNamespace,
			"Namespace":            g.appConfig.DbAuthNamespace,
			"AuthNamespace":        g.appConfig.DbAuthNamespace,
			"Username":             g.appConfig.DbUsername,
		},
	}, embedMigrations)
	if err != nil {
		return err
	}

	g.log.Info("auth migrations applied successfully")

	return nil
}

func (g *Auth) Handler() http.Handler {
	return g.handler
}
