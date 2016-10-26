package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

type apiServer struct {
	system        *gohome.System
	recipeManager *gohome.RecipeManager
	eventLogger   gohome.WSEventLogger
}

// ListenAndServe creates a new WWW server, that handles API calls and also
// runs the gohome website
func ListenAndServe(
	port string,
	system *gohome.System,
	recipeManager *gohome.RecipeManager,
	eventLogger gohome.WSEventLogger) error {
	server := &apiServer{
		system:        system,
		recipeManager: recipeManager,
		eventLogger:   eventLogger,
	}
	return server.listenAndServe(port)
}

func (s *apiServer) listenAndServe(port string) error {

	r := mux.NewRouter()

	//TODO: re-enable
	//r.HandleFunc("/api/v1/events/ws", s.eventLogger.HTTPHandler())

	RegisterSceneHandlers(r, s)
	RegisterButtonHandlers(r, s)
	RegisterZoneHandlers(r, s)
	RegisterDeviceHandlers(r, s)
	RegisterSensorHandlers(r, s)
	RegisterDiscoveryHandlers(r, s)
	RegisterCookBookHandlers(r, s)
	RegisterRecipeHandlers(r, s)
	RegisterMonitorHandlers(r, s)

	return http.ListenAndServe(
		port,
		handlers.CORS(
			handlers.AllowedMethods([]string{"PUT", "POST", "DELETE", "GET", "OPTIONS", "UPGRADE"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(r))
}