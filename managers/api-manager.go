package managers

import (
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
)

func ApiManagerGet(apiName string) (models.ApiModel, string, error) {
	var a models.ApiModel

	aDb, message, err := repositories.ApiGet(apiName)
	if err != nil {
		return a, message, err
	}

	return aDb, message, err
}
