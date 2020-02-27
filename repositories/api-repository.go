package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func apiExists(a models.ApiModel) bool {
	var aDb models.ApiModel
	api.Api.Database.Orm.Where(&models.ApiModel{
		Name: a.Name,
	}).Find(&aDb)

	return aDb.Name == a.Name
}

func ApiGet(apiName string) (models.ApiModel, string, error) {
	var a models.ApiModel
	api.Api.Database.Orm.Where(&models.ApiModel{
		Name: apiName,
	}).First(&a)

	if !apiExists(a) {
		msg := fmt.Sprintf("no API with name %s found", apiName)
		return a, msg, errors.New(msg)
	}

	return a, "api retrieved", nil
}
