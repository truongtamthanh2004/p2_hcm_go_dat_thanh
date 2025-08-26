package routes

import (
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Any("/api/v1/auth/*path", proxy.NewReverseProxy(os.Getenv("AUTH_SERVICE_URL")))

	r.Any("/api/v1/users/*path",
		proxy.NewReverseProxy(os.Getenv("USER_SERVICE_URL")),
	)

	r.Any("/api/venues", proxy.NewReverseProxy(os.Getenv("VENUE_SERVICE_URL")))
	r.Any("/api/venues/*path", proxy.NewReverseProxy(os.Getenv("VENUE_SERVICE_URL")))

	r.Any("/api/booking/*path", middleware.AuthMiddleware(), proxy.NewReverseProxy(os.Getenv("BOOKING_SERVICE_URL")))

	r.Any("/api/payment/*path", middleware.AuthMiddleware(), proxy.NewReverseProxy(os.Getenv("PAYMENT_SERVICE_URL")))

	r.Any("/api/v1/chat/*path", middleware.AuthMiddleware(), proxy.NewReverseProxy(os.Getenv("CHAT_SERVICE_URL")))

	r.Any("/api/notify/*path", middleware.AuthMiddleware(), proxy.NewReverseProxy(os.Getenv("NOTIFICATION_SERVICE_URL")))

	r.Any("/api/map/*path", proxy.NewReverseProxy(os.Getenv("MAP_SERVICE_URL")))

	r.Any("/api/stats/*path",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin"),
		proxy.NewReverseProxy(os.Getenv("STATISTIC_SERVICE_URL")),
	)
}
