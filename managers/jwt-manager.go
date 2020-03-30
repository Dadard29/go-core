package managers

import (
	"errors"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
)

func GetJwtSecret() string {
	jwtSecretKey := "JWT_SECRET"

	secret := api.Api.Config.GetEnv(jwtSecretKey)
	if secret == "" {
		logger.CheckErrFatal(errors.New("no configured jwt secret"))
	}

	return secret
}

func ValidateJwtCiphered(jwtCiphered string) *auth.JwtPayload {
	jwtDeciphered, err := auth.DecipherJwtWithJwe(config.PrivateKeyFile, []byte(jwtCiphered))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	secret := GetJwtSecret()
	p, err := auth.VerifyJwtHS256(jwtDeciphered, secret)
	if err != nil {
		return nil
	}

	return p
}
