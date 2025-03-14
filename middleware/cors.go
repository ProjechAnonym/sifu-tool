package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
)

func Cors(domains []string) cors.Config {
	config := cors.Config{
		AllowOrigins:     domains,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "domain", "scheme", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return config
}