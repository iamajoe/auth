package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamajoe/auth/internal/conf"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type GlobalConfiguration = conf.GlobalConfiguration
type APIConfiguration = conf.APIConfiguration
type DBConfiguration = conf.DBConfiguration
type ProviderConfiguration = conf.ProviderConfiguration
type LoggingConfig = conf.LoggingConfig
type ProfilerConfig = conf.ProfilerConfig
type TracingConfig = conf.TracingConfig
type MetricsConfig = conf.MetricsConfig
type SMTPConfiguration = conf.SMTPConfiguration
type PasswordConfiguration = conf.PasswordConfiguration
type JWTConfiguration = conf.JWTConfiguration
type MailerConfiguration = conf.MailerConfiguration
type SmsProviderConfiguration = conf.SmsProviderConfiguration
type HookConfiguration = conf.HookConfiguration
type SecurityConfiguration = conf.SecurityConfiguration
type SessionsConfiguration = conf.SessionsConfiguration
type MFAConfiguration = conf.MFAConfiguration
type SAMLConfiguration = conf.SAMLConfiguration
type CORSConfiguration = conf.CORSConfiguration

type JwtKeysDecoder = conf.JwtKeysDecoder
type JwkInfo = conf.JwkInfo

func GetSigningAlg(k jwk.Key) jwt.SigningMethod {
	return conf.GetSigningAlg(k)
}
