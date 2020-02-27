package managers

import (
	"errors"
	"github.com/Dadard29/go-api-utils/auth"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"net/http"
)

func GetJwtSecret() string {
	jwtSecretKey, err := api.Api.Config.GetValueFromFile(
		config.Profile,
		config.ProfileJwt,
		config.ProfileJwtSecretKey)

	logger.CheckErrFatal(err)

	secret := api.Api.Config.GetEnv(jwtSecretKey)
	if secret == "" {
		logger.CheckErrFatal(errors.New("no configured jwt secret"))
	}

	return secret
}

func ValidateJwtBody(jwtCiphered string) (bool, int, string) {
	jwtDeciphered, err := auth.DecipherJwtWithJwe(config.PrivateKeyFile, []byte(jwtCiphered))
	if err != nil {
		logger.Error(err.Error())
		return false, http.StatusInternalServerError, "error deciphering token"
	}

	secret := GetJwtSecret()
	_, err = auth.VerifyJwtHS256(jwtDeciphered, secret)
	if err != nil {
		return false, http.StatusBadRequest, "token invalid"
	}

	return true, http.StatusOK, "token valid"
}
