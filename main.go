package main

import (
	"github.com/Dadard29/go-api-utils/API"
	"github.com/Dadard29/go-api-utils/service"
	. "github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/controllers"
	"net/http"
)

//var routes = map[string]func(w http.ResponseWriter, r *http.Request) {
var routes = service.RouteMapping{
	Mapping: map[string]service.Route{
		"/notification/bot/webhook": {Handler: controllers.NotificationBotWebookRoute, Method:  http.MethodGet},
	},
}


func main() {

	Api = API.NewAPI("Core", "config/config.json", routes, true)

	Api.Service.Start()

	Api.Service.Stop()
}
