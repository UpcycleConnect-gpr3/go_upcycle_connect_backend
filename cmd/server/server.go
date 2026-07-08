package server

import (
	"go-upcycle_connect-backend/app/handlers/appointment_handlers"
	"go-upcycle_connect-backend/app/handlers/stats_handlers"
	"go-upcycle_connect-backend/app/handlers/activity_handlers"
	"go-upcycle_connect-backend/app/handlers/ad_handlers"
	"go-upcycle_connect-backend/app/handlers/auth_handlers"
	"go-upcycle_connect-backend/app/handlers/billing_handlers"
	"go-upcycle_connect-backend/app/handlers/delivery_method_handlers"
	"go-upcycle_connect-backend/app/handlers/event_handlers"
	"go-upcycle_connect-backend/app/handlers/event_step_handlers"
	"go-upcycle_connect-backend/app/handlers/invoice_handlers"
	"go-upcycle_connect-backend/app/handlers/locker_handlers"
	"go-upcycle_connect-backend/app/handlers/metric_handlers"
	"go-upcycle_connect-backend/app/handlers/object_handlers"
	"go-upcycle_connect-backend/app/handlers/object_order_handlers"
	"go-upcycle_connect-backend/app/handlers/order_delivery_method_handlers"
	"go-upcycle_connect-backend/app/handlers/notification_handlers"
	"go-upcycle_connect-backend/app/handlers/order_handlers"
	"go-upcycle_connect-backend/app/handlers/package_handlers"
	"go-upcycle_connect-backend/app/handlers/prestataire_handlers"
	"go-upcycle_connect-backend/app/handlers/payment_handlers"
	"go-upcycle_connect-backend/app/handlers/project_handlers"
	"go-upcycle_connect-backend/app/handlers/score_handlers"
	"go-upcycle_connect-backend/app/handlers/step_handlers"
	"go-upcycle_connect-backend/app/handlers/upload_handlers"
	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/middleware/ratelimit_middleware"
	"go-upcycle_connect-backend/config"
	"go-upcycle_connect-backend/var/database"
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
	config.InitEmail()
	config.InitFilesystem()

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

	http.HandleFunc("POST /auth/login/{$}", limiterMedium.RateLimit(auth_handlers.LoginHandler))
	http.HandleFunc("PATCH /user/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(auth_handlers.UpdateUserHandler)))

	// Delivery Method routes
	http.HandleFunc("GET /delivery-methods", limiterHigh.RateLimit(delivery_method_handlers.IndexDeliveryMethodHandler))
	http.HandleFunc("GET /delivery-methods/{id}", limiterHigh.RateLimit(delivery_method_handlers.ShowDeliveryMethodHandler))
	http.HandleFunc("POST /delivery-methods", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.StoreDeliveryMethodHandler)))
	http.HandleFunc("PUT /delivery-methods/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.UpdateDeliveryMethodHandler)))
	http.HandleFunc("DELETE /delivery-methods/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(delivery_method_handlers.DeleteDeliveryMethodHandler)))

	// Event routes
	http.HandleFunc("GET /events", limiterHigh.RateLimit(event_handlers.IndexEventHandler))
	http.HandleFunc("GET /events/{id}", limiterHigh.RateLimit(event_handlers.ShowEventHandler))
	http.HandleFunc("POST /events", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.StoreEventHandler)))
	http.HandleFunc("PUT /events/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.UpdateEventHandler)))
	http.HandleFunc("DELETE /events/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.DeleteEventHandler)))
	http.HandleFunc("GET /events/{id}/steps", limiterHigh.RateLimit(event_handlers.GetEventStepsHandler))
	http.HandleFunc("POST /events/{id}/steps", limiterMedium.RateLimit(auth_middleware.IsAuth(event_handlers.CreateEventStepHandler)))

	// Event Step routes
	http.HandleFunc("GET /event-steps", limiterHigh.RateLimit(event_step_handlers.IndexEventStepHandler))
	http.HandleFunc("GET /event-steps/{id}", limiterHigh.RateLimit(event_step_handlers.ShowEventStepHandler))
	http.HandleFunc("POST /event-steps", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.StoreEventStepHandler)))
	http.HandleFunc("PUT /event-steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.UpdateEventStepHandler)))
	http.HandleFunc("DELETE /event-steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(event_step_handlers.DeleteEventStepHandler)))

	// Object routes
	http.HandleFunc("GET /objects", limiterHigh.RateLimit(object_handlers.IndexObjectHandler))
	http.HandleFunc("GET /objects/{id}", limiterHigh.RateLimit(object_handlers.ShowObjectHandler))
	http.HandleFunc("POST /objects", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.StoreObjectHandler)))
	http.HandleFunc("PUT /objects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.UpdateObjectHandler)))
	http.HandleFunc("DELETE /objects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.DeleteObjectHandler)))
	http.HandleFunc("POST /objects/{id}/validate", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.ValidateObjectHandler)))
	http.HandleFunc("POST /objects/{id}/reject", limiterMedium.RateLimit(auth_middleware.IsAuth(object_handlers.RejectObjectHandler)))
	http.HandleFunc("GET /score/config", limiterHigh.RateLimit(score_handlers.GetScoreConfigHandler))
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
	http.HandleFunc("GET /projects", limiterHigh.RateLimit(project_handlers.IndexProjectHandler))
	http.HandleFunc("GET /projects/{id}", limiterHigh.RateLimit(project_handlers.ShowProjectHandler))
	http.HandleFunc("POST /projects", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.StoreProjectHandler)))
	http.HandleFunc("PUT /projects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.UpdateProjectHandler)))
	http.HandleFunc("DELETE /projects/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.DeleteProjectHandler)))
	http.HandleFunc("GET /projects/{id}/objects", limiterHigh.RateLimit(project_handlers.GetProjectObjectsHandler))
	http.HandleFunc("POST /projects/{id}/objects/{objectId}", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.LinkObjectHandler)))
	http.HandleFunc("GET /projects/{id}/steps", limiterHigh.RateLimit(project_handlers.GetProjectStepsHandler))
	http.HandleFunc("POST /projects/{id}/steps", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.CreateProjectStepHandler)))
	http.HandleFunc("POST /projects/{id}/feature", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.FeatureProjectHandler)))
	http.HandleFunc("DELETE /projects/{id}/feature", limiterMedium.RateLimit(auth_middleware.IsAuth(project_handlers.UnfeatureProjectHandler)))

	// Step routes
	http.HandleFunc("GET /steps", limiterHigh.RateLimit(step_handlers.IndexStepHandler))
	http.HandleFunc("GET /steps/{id}", limiterHigh.RateLimit(step_handlers.ShowStepHandler))
	http.HandleFunc("POST /steps", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.StoreStepHandler)))
	http.HandleFunc("PUT /steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.UpdateStepHandler)))
	http.HandleFunc("DELETE /steps/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(step_handlers.DeleteStepHandler)))

	// Ads routes (publicites du professionnel connecte)
	http.HandleFunc("GET /ads", limiterMedium.RateLimit(auth_middleware.IsAuth(ad_handlers.IndexAdHandler)))
	http.HandleFunc("POST /ads", limiterMedium.RateLimit(auth_middleware.IsAuth(ad_handlers.StoreAdHandler)))
	http.HandleFunc("PUT /ads/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(ad_handlers.UpdateAdHandler)))
	http.HandleFunc("DELETE /ads/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(ad_handlers.DeleteAdHandler)))

	// Stats & abonnement de l'utilisateur connecte
	http.HandleFunc("GET /stats/me", limiterHigh.RateLimit(auth_middleware.IsAuth(stats_handlers.GetMyStatsHandler)))
	http.HandleFunc("GET /stats/finance", limiterHigh.RateLimit(auth_middleware.IsAuth(stats_handlers.GetFinanceStatsHandler)))
	http.HandleFunc("GET /subscriptions/me", limiterHigh.RateLimit(auth_middleware.IsAuth(billing_handlers.GetMySubscriptionHandler)))
	http.HandleFunc("GET /invoices/me", limiterHigh.RateLimit(auth_middleware.IsAuth(invoice_handlers.IndexInvoicesHandler)))
	http.HandleFunc("GET /invoices/{ref}/pdf", limiterHigh.RateLimit(auth_middleware.IsAuth(invoice_handlers.InvoicePdfHandler)))

	// Appointment routes (planning personnel de l'utilisateur du token)
	http.HandleFunc("GET /appointments", limiterMedium.RateLimit(auth_middleware.IsAuth(appointment_handlers.IndexAppointmentHandler)))
	http.HandleFunc("POST /appointments", limiterMedium.RateLimit(auth_middleware.IsAuth(appointment_handlers.StoreAppointmentHandler)))
	http.HandleFunc("PUT /appointments/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(appointment_handlers.UpdateAppointmentHandler)))
	http.HandleFunc("DELETE /appointments/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(appointment_handlers.DeleteAppointmentHandler)))

	// Upload routes (images d'annonces)
	http.HandleFunc("POST /upload", limiterMedium.RateLimit(auth_middleware.IsAuth(upload_handlers.StoreUploadHandler)))
	http.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Order routes
	http.HandleFunc("GET /orders/me", limiterHigh.RateLimit(auth_middleware.IsAuth(order_handlers.MyOrdersHandler)))
	http.HandleFunc("GET /orders", limiterHigh.RateLimit(order_handlers.IndexOrderHandler))
	http.HandleFunc("GET /orders/{id}", limiterHigh.RateLimit(order_handlers.ShowOrderHandler))
	http.HandleFunc("POST /orders", limiterMedium.RateLimit(auth_middleware.IsAuth(order_handlers.StoreOrderHandler)))
	http.HandleFunc("PUT /orders/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(order_handlers.UpdateOrderHandler)))
	http.HandleFunc("DELETE /orders/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(order_handlers.DeleteOrderHandler)))

	// Prestataire routes
	http.HandleFunc("GET /prestataires", limiterHigh.RateLimit(prestataire_handlers.IndexPrestataireHandler))
	http.HandleFunc("GET /prestataires/{id}", limiterHigh.RateLimit(prestataire_handlers.ShowPrestataireHandler))
	http.HandleFunc("POST /prestataires", limiterMedium.RateLimit(auth_middleware.IsAuth(prestataire_handlers.StorePrestataireHandler)))
	http.HandleFunc("PUT /prestataires/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(prestataire_handlers.UpdatePrestataireHandler)))
	http.HandleFunc("DELETE /prestataires/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(prestataire_handlers.DeletePrestataireHandler)))

	// Notification routes
	http.HandleFunc("GET /notifications/me", limiterHigh.RateLimit(auth_middleware.IsAuth(notification_handlers.MyNotificationsHandler)))
	http.HandleFunc("POST /notifications", limiterMedium.RateLimit(auth_middleware.IsAuth(notification_handlers.StoreNotificationHandler)))
	http.HandleFunc("POST /notifications/{id}/read", limiterMedium.RateLimit(auth_middleware.IsAuth(notification_handlers.MarkNotificationReadHandler)))

	// Activity log routes
	http.HandleFunc("GET /logs", limiterHigh.RateLimit(auth_middleware.IsAuth(activity_handlers.GetLogsHandler)))

	// Locker routes
	http.HandleFunc("GET /lockers", limiterHigh.RateLimit(locker_handlers.IndexLockerHandler))
	http.HandleFunc("GET /lockers/{id}", limiterHigh.RateLimit(locker_handlers.ShowLockerHandler))
	http.HandleFunc("POST /lockers", limiterMedium.RateLimit(auth_middleware.IsAuth(locker_handlers.StoreLockerHandler)))
	http.HandleFunc("PUT /lockers/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(locker_handlers.UpdateLockerHandler)))
	http.HandleFunc("DELETE /lockers/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(locker_handlers.DeleteLockerHandler)))

	// Package routes
	http.HandleFunc("GET /packages", limiterHigh.RateLimit(package_handlers.IndexPackageHandler))
	http.HandleFunc("GET /packages/{id}", limiterHigh.RateLimit(package_handlers.ShowPackageHandler))
	http.HandleFunc("POST /packages", limiterMedium.RateLimit(auth_middleware.IsAuth(package_handlers.StorePackageHandler)))
	http.HandleFunc("PUT /packages/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(package_handlers.UpdatePackageHandler)))
	http.HandleFunc("DELETE /packages/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(package_handlers.DeletePackageHandler)))

	// Object-Order routes
	http.HandleFunc("GET /object-orders", limiterHigh.RateLimit(object_order_handlers.IndexObjectOrderHandler))
	http.HandleFunc("GET /object-orders/{id}", limiterHigh.RateLimit(object_order_handlers.ShowObjectOrderHandler))
	http.HandleFunc("POST /object-orders", limiterMedium.RateLimit(auth_middleware.IsAuth(object_order_handlers.StoreObjectOrderHandler)))
	http.HandleFunc("DELETE /object-orders/{id}", limiterMedium.RateLimit(auth_middleware.IsAuth(object_order_handlers.DeleteObjectOrderHandler)))

	// Order-DeliveryMethod routes
	http.HandleFunc("GET /order-delivery-methods", limiterHigh.RateLimit(order_delivery_method_handlers.IndexOrderDeliveryMethodHandler))
	http.HandleFunc("POST /order-delivery-methods", limiterMedium.RateLimit(auth_middleware.IsAuth(order_delivery_method_handlers.StoreOrderDeliveryMethodHandler)))
	http.HandleFunc("DELETE /order-delivery-methods/{orderId}/{deliveryMethodId}", limiterMedium.RateLimit(auth_middleware.IsAuth(order_delivery_method_handlers.DeleteOrderDeliveryMethodHandler)))

	// Billing / Stripe routes
	http.HandleFunc("POST /billing/checkout-session", limiterMedium.RateLimit(auth_middleware.IsAuth(billing_handlers.CreateCheckoutSessionHandler)))
	http.HandleFunc("GET /billing/checkout-session/{id}", limiterHigh.RateLimit(auth_middleware.IsAuth(billing_handlers.GetCheckoutSessionHandler)))
	// Payment routes (paiement unique d'une annonce)
	http.HandleFunc("POST /payments/checkout", limiterMedium.RateLimit(auth_middleware.IsAuth(payment_handlers.CreatePaymentHandler)))
	http.HandleFunc("GET /payments/session/{id}", limiterHigh.RateLimit(auth_middleware.IsAuth(payment_handlers.GetPaymentStatusHandler)))

	// Webhook: no auth middleware (Stripe sends no JWT) — verified by signature.
	// Single endpoint for both subscriptions and one-time annonce payments.
	http.HandleFunc("POST /billing/webhook", billing_handlers.WebhookHandler)

	logger.Info().Msg("Listening at http://localhost:" + os.Getenv("APP_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), corsMiddleware(http.DefaultServeMux))
	if err != nil {
		return
	}
}

// corsMiddleware enables cross-origin requests from the local dev frontends
// (Vite on a different port). Reflects the request Origin and answers the
// preflight OPTIONS so browser calls are not blocked.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
