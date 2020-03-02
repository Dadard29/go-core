package managers

import (
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/models"
	"time"
)

func SessionManagerCheckAuth(username string, password string) (bool, error) {

	usernameKey, err := api.Api.Config.GetValueFromFile(
		config.Session,
		config.SessionAuth,
		config.SessionAuthUsernameKey)
	if err != nil {
		return false, err
	}

	passwordKey, err := api.Api.Config.GetValueFromFile(
		config.Session,
		config.SessionAuth,
		config.SessionAuthUPasswordKey)
	if err != nil {
		return false, err
	}

	usernameEnv := api.Api.Config.GetEnv(usernameKey)
	passwordEnv := api.Api.Config.GetEnv(passwordKey)

	if usernameEnv == username && passwordEnv == password {
		return true, nil
	} else {
		return false, nil
	}
}

func SessionManagerGet() (bool, string, error) {
	return false, "", nil
}

func SessionManagerCreate(duration time.Duration) (bool, string, error) {
	var s models.Session
	s.AccessToken = s.NewAccessToken()
	s.CreatedAt = time.Now()
	s.Duration = duration

	//sessionDb, msg, err := repositories.
	return false, "", nil
}
