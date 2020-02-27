package managers

import (
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
)

func ApiManagerList() ([]models.ApiModel, string, error) {
	var aList []models.ApiModel

	aDb, message, err := repositories.ApiGetList()
	if err != nil {
		return aList, message, err
	}

	return aDb, message, nil
}

func ApiManagerGet(apiName string) (models.ApiModel, string, error) {
	var a models.ApiModel

	aDb, message, err := repositories.ApiGet(apiName)
	if err != nil {
		return a, message, err
	}

	return aDb, message, nil
}
