package repositories

import (
	"errors"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func sessionExistsFromExpiration() bool {
	var sessions []models.Session
	api.Api.Database.Orm.Find(&sessions)

	//for _, s := range sessions {
	//
	//}
	return false
}

func sessionExistsFromToken(s models.Session) bool {
	var sessionDb models.Session
	api.Api.Database.Orm.Where(&models.Session{
		AccessToken: s.AccessToken,
	}).Find(&sessionDb)

	return sessionDb.AccessToken == s.AccessToken
}

func SessionCreate(s models.Session) (models.Session, string, error) {
	if sessionExistsFromToken(s) {
		msg := "existing session"
		return models.Session{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Create(&s)

	if sessionExistsFromToken(s) {
		msg := "error creating session"
		return models.Session{}, msg, errors.New(msg)
	}

	return s, "session created", nil
}
