package client

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/iamajoe/auth/pkg"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

func GenerateJWTKeys(secret string, keyID string) (pkg.JwtKeysDecoder, []string, error) {
	if secret == "" {
		return nil, nil, fmt.Errorf("secret cannot be empty")
	}

	keys := make(pkg.JwtKeysDecoder)
	validMethods := []string{}

	// transform the secret into a JWK for consistency
	privKey, err := jwk.FromRaw([]byte(secret))
	if err != nil {
		return keys, validMethods, err
	}

	if keyID != "" {
		if err := privKey.Set(jwk.KeyIDKey, keyID); err != nil {
			return keys, validMethods, err
		}
	}

	if privKey.Algorithm().String() == "" {
		if err := privKey.Set(jwk.AlgorithmKey, jwt.SigningMethodHS256.Name); err != nil {
			return keys, validMethods, err
		}
	}

	if err := privKey.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return keys, validMethods, err
	}

	if len(privKey.KeyOps()) == 0 {
		if err := privKey.Set(jwk.KeyOpsKey, jwk.KeyOperationList{jwk.KeyOpSign, jwk.KeyOpVerify}); err != nil {
			return keys, validMethods, err
		}
	}

	pubKey, err := privKey.PublicKey()
	if err != nil {
		return keys, validMethods, err
	}

	keys[keyID] = pkg.JwkInfo{
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}

	for _, key := range keys {
		alg := pkg.GetSigningAlg(key.PublicKey)
		validMethods = append(validMethods, alg.Alg())
	}

	return keys, validMethods, nil
}
