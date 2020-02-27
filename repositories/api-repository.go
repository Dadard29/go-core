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

func ApiCreate(a models.ApiModel) (models.ApiModel, string, error) {
	if apiExists(a) {
		msg := "existing api with same name"
		return models.ApiModel{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Create(&a)

	if !apiExists(a) {
		msg := "failed to create api"
		return models.ApiModel{}, msg, errors.New(msg)
	}

	return a, "api created", nil
}

func ApiUpdate(a models.ApiModel) (models.ApiModel, string, error) {
	if !apiExists(a) {
		msg := fmt.Sprintf("no api with name %s found", a.Name)
		return models.ApiModel{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Save(&a)

	return a, "api updated", nil
}

func ApiDelete(a models.ApiModel) (models.ApiModel, string, error) {
	if !apiExists(a) {
		msg := fmt.Sprintf("no api with name %s found", a.Name)
		return models.ApiModel{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Delete(&a)

	return a, "api deleted", nil
}
