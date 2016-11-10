package pinger

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route


var routes = Routes{
	Route{
		"GET_PINGS",
		"GET",
		"/api/pings",
		Pings,
	},
	Route{
		"CFG.GET.MACHINES",
		"GET",
		"/api/cfg/machines",
		CFGGetMachines,
	},
	Route{
		"CFG.ADD.MACHINES",
		"PUT",
		"/api/cfg/machines/{name}",
		CFGAddMachine,
	},
	Route{
		"CFG.GET",
		"GET",
		"/api/cfg",
		CFGGet,
	},
}
