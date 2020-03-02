package main

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/database"
	"github.com/Dadard29/go-api-utils/service"
	. "github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/controllers"
	"github.com/Dadard29/go-core/models"
	"net/http"
)

//var routes = map[string]func(w http.ResponseWriter, r *http.Request) {
var routes = service.RouteMapping{
	Mapping: map[string]service.Route{
		"/notification/bot/webhook": {
			Handler: controllers.NotificationBotWebookRoute,
			Method:  []string{http.MethodPost},
		},
		"/profile/auth/jwt": {
			Handler: controllers.JwtHandler,
			Method:  []string{http.MethodGet, http.MethodPost},
		},
		"/profile/auth": {
			Handler: controllers.ProfileHandler,
			Method:  []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		},
		"/api": {
			Handler: controllers.ApiHandler,
			Method:  []string{http.MethodGet},
		},
		"/subs": {
			Handler: controllers.SubsHandler,
			Method:  []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		},
		"/session": {
			Handler: controllers.SessionHandler,
			Method:  []string{http.MethodGet, http.MethodPost},
		},
	},
}

func main() {

	Api = API.NewAPI("Core", "config/config.json", routes, true)

	dbConfig, err := Api.Config.GetSubcategoryFromFile("api", "db")
	Api.Logger.CheckErr(err)
	Api.Database = database.NewConnector(dbConfig, true, []interface{}{
		models.Profile{},
		models.ApiModel{},
		models.Subscription{},
	})

	Api.Service.Start()

	Api.Service.Stop()
}
