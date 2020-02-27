package managers

import (
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
)

func SubsManagerList(profileKey string) ([]models.Subscription, string, error) {
	var s []models.Subscription

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	subDb, msg, err := repositories.SubsListFromProfile(p)
	if err != nil {
		return s, msg, err
	}

	return subDb, msg, nil
}
