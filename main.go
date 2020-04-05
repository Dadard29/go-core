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
	"/notification/bot/webhook": {
		Description: "manage notifications from webhooks",
		MethodMapping: service.MethodMapping{
			http.MethodPost: controllers.NotificationBotWebookRoute,
		},
	},
	"/profile/auth/jwt": {
		Description: "manage the web tokens",
		MethodMapping: service.MethodMapping{
			http.MethodPost: controllers.JwtCreate,
			http.MethodGet:  controllers.JwtValidate,
		},
	},
	"/profile/auth": {
		Description: "manage the profile DB object",
		MethodMapping: service.MethodMapping{
			http.MethodPost:   controllers.ProfileSignUp,
			http.MethodGet:    controllers.ProfileGet,
			http.MethodPut:    controllers.ProfileChangePassword,
			http.MethodDelete: controllers.ProfileDelete,
		},
	},
	"/profile/auth/confirm": {
		Description:   "manage confirmation for account creation",
		MethodMapping: service.MethodMapping{
			http.MethodPost: controllers.ProfileSignUpConfirm,
		},
	},
	"/api": {
		Description: "manage the APIs DB objects",
		MethodMapping: service.MethodMapping{
			http.MethodGet: controllers.ApiGet,
		},
	},
	"/api/list": {
		Description: "list the APIs DB objects",
		MethodMapping: service.MethodMapping{
			http.MethodGet: controllers.ApiListGet,
		},
	},
	"/subs": {
		Description: "manage the subscriptions DB objects",
		MethodMapping: service.MethodMapping{
			http.MethodPost:   controllers.Subscribe,
			http.MethodGet:    controllers.SubsCheckExists,
			http.MethodPut: controllers.SubRegenerate,
		},
	},
	"/subs/list": {
		Description: "manage the subscriptions list",
		MethodMapping: service.MethodMapping{
			http.MethodGet: controllers.SubsList,
		},
	},
}

// ENV

// BOT_WEBOOK: 		secret for authenticating gitlab webhook
// CI_BOT_TOKEN: 	token of bot dedicated for CI
// STD_BOT_TOKEN: 	token of bot dedicated for dadard.website notifs
// USERNAME_DB: 	username for database
// PASSWORD_DB: 	password for database
// JWT_SECRET: 		secret to generate jwt
// VERSION: 		version
// CORS_ORIGIN: 	specify origin for web access
// SMTP_PASSWORD: 	password for sending mails

func main() {

	Api = API.NewAPI("Core", "config/config.json", routes, true)

	dbConfig, err := Api.Config.GetSubcategoryFromFile("api", "db")
	Api.Logger.CheckErr(err)
	Api.Database = database.NewConnector(dbConfig, true, []interface{}{
		models.Profile{},
		models.TempProfile{},
		models.ApiModel{},
		models.Subscription{},
	})

	Api.Service.Start()

	Api.Service.Stop()
}
