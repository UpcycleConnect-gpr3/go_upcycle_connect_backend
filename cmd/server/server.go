package server

import (
	"go-upcycle_connect-backend/app/handlers/delivery_method_handlers"
	"go-upcycle_connect-backend/app/handlers/event_handlers"
	"go-upcycle_connect-backend/app/handlers/event_step_handlers"
	"go-upcycle_connect-backend/app/handlers/metric_handlers"
	"go-upcycle_connect-backend/app/handlers/object_handlers"
	"go-upcycle_connect-backend/app/handlers/project_handlers"
	"go-upcycle_connect-backend/app/handlers/step_handlers"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/middleware/ratelimit_middleware"
	"go-upcycle_connect-backend/config"
	"go-upcycle_connect-backend/database"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	log "github.com/thedataflows/go-lib-log"
)

func initialize() {

	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithLogFormat(log.LOG_FORMAT_JSON).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Config Initialization
	config.InitDatabase()

	err = database.UpcycleConnect.Ping()

	if err != nil {
		logger.Fatal().Err(err).Msg("(DATABASE)")
	}
}

func Start() {

	initialize()

	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	// Rate limiting
	limiterMedium := ratelimit_middleware.NewRateLimiter(30, 1*time.Minute)
	limiterHigh := ratelimit_middleware.NewRateLimiter(60, 1*time.Minute)

	// Container
	//containerBackoffice := source_middleware.Container("go-upcycle_connect-backend")

	// Health
	http.HandleFunc("GET /health/{$}", metric_handlers.Health)

	// Delivery Method routes
	http.HandleFunc("GET /delivery-methods", limiterHigh.RateLimit(delivery_method_handlers.GetDeliveryMethodsHandler))
	http.HandleFunc("GET /delivery-methods/{id}", limiterHigh.RateLimit(delivery_method_handlers.GetDeliveryMethodHandler))
	http.HandleFunc("POST /delivery-methods", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.CreateDeliveryMethodHandler)))
	http.HandleFunc("PUT /delivery-methods/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.UpdateDeliveryMethodHandler)))
	http.HandleFunc("DELETE /delivery-methods/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.DeleteDeliveryMethodHandler)))

	// Event routes
	http.HandleFunc("GET /events", limiterHigh.RateLimit(event_handlers.GetEventsHandler))
	http.HandleFunc("GET /events/{id}", limiterHigh.RateLimit(event_handlers.GetEventHandler))
	http.HandleFunc("POST /events", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.CreateEventHandler)))
	http.HandleFunc("PUT /events/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.UpdateEventHandler)))
	http.HandleFunc("DELETE /events/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.DeleteEventHandler)))
	http.HandleFunc("GET /events/{id}/steps", limiterHigh.RateLimit(event_handlers.GetEventStepsHandler))
	http.HandleFunc("POST /events/{id}/steps", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.CreateEventStepHandler)))

	// Event Step routes
	http.HandleFunc("GET /event-steps", limiterHigh.RateLimit(event_step_handlers.GetEventStepsHandler))
	http.HandleFunc("GET /event-steps/{id}", limiterHigh.RateLimit(event_step_handlers.GetEventStepHandler))
	http.HandleFunc("POST /event-steps", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.CreateEventStepHandler)))
	http.HandleFunc("PUT /event-steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.UpdateEventStepHandler)))
	http.HandleFunc("DELETE /event-steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.DeleteEventStepHandler)))

	// Object routes
	http.HandleFunc("GET /objects", limiterHigh.RateLimit(object_handlers.GetObjectsHandler))
	http.HandleFunc("GET /objects/{id}", limiterHigh.RateLimit(object_handlers.GetObjectHandler))
	http.HandleFunc("POST /objects", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.CreateObjectHandler)))
	http.HandleFunc("PUT /objects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.UpdateObjectHandler)))
	http.HandleFunc("DELETE /objects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.DeleteObjectHandler)))
	http.HandleFunc("GET /objects/{id}/score", limiterHigh.RateLimit(object_handlers.GetObjectScoreHandler))
	http.HandleFunc("GET /objects/{id}/delivery-methods", limiterHigh.RateLimit(object_handlers.GetObjectDeliveryMethodsHandler))
	http.HandleFunc("POST /objects/{id}/delivery-methods/{deliveryMethodId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.LinkDeliveryMethodHandler)))
	http.HandleFunc("DELETE /objects/{id}/delivery-methods/{deliveryMethodId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.UnlinkDeliveryMethodHandler)))
	http.HandleFunc("GET /objects/{id}/projects", limiterHigh.RateLimit(object_handlers.GetObjectProjectsHandler))
	http.HandleFunc("POST /objects/{id}/projects/{projectId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.LinkProjectHandler)))
	http.HandleFunc("DELETE /objects/{id}/projects/{projectId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.UnlinkProjectHandler)))
	http.HandleFunc("GET /objects/{id}/users", limiterHigh.RateLimit(object_handlers.GetObjectUsersHandler))
	http.HandleFunc("POST /objects/{id}/users/{userId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.LinkUserHandler)))
	http.HandleFunc("DELETE /objects/{id}/users/{userId}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.UnlinkUserHandler)))

	// Project routes
	http.HandleFunc("GET /projects", limiterHigh.RateLimit(project_handlers.GetProjectsHandler))
	http.HandleFunc("GET /projects/{id}", limiterHigh.RateLimit(project_handlers.GetProjectHandler))
	http.HandleFunc("POST /projects", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.CreateProjectHandler)))
	http.HandleFunc("PUT /projects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.UpdateProjectHandler)))
	http.HandleFunc("DELETE /projects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.DeleteProjectHandler)))
	http.HandleFunc("GET /projects/{id}/score", limiterHigh.RateLimit(project_handlers.GetProjectScoreHandler))
	http.HandleFunc("GET /projects/{id}/objects", limiterHigh.RateLimit(project_handlers.GetProjectObjectsHandler))
	http.HandleFunc("POST /projects/{id}/objects/{objectId}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.LinkObjectHandler)))
	http.HandleFunc("GET /projects/{id}/steps", limiterHigh.RateLimit(project_handlers.GetProjectStepsHandler))
	http.HandleFunc("POST /projects/{id}/steps", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.CreateProjectStepHandler)))

	// Step routes
	http.HandleFunc("GET /steps", limiterHigh.RateLimit(step_handlers.GetStepsHandler))
	http.HandleFunc("GET /steps/{id}", limiterHigh.RateLimit(step_handlers.GetStepHandler))
	http.HandleFunc("POST /steps", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.CreateStepHandler)))
	http.HandleFunc("PUT /steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.UpdateStepHandler)))
	http.HandleFunc("DELETE /steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.DeleteStepHandler)))

	logger.Info().Msg("Listening at http://localhost:" + os.Getenv("APP_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
	if err != nil {
		return
	}
}
