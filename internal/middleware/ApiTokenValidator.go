package middleware

import (
	"net/http"

	"github.com/handrixn/task-tracker/internal/util"
	"github.com/spf13/viper"
)

func ValidateAPIToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiToken := r.Header.Get("x-api-token")

		if apiToken == "" {
			apiToken = r.Header.Get("X-API-TOKEN")

			if apiToken == "" {
				util.JsonResponse(w, http.StatusUnauthorized, nil)
				return
			}
		}

		apiTokenSource := viper.GetString("API_TOKEN")
		if apiToken != apiTokenSource {
			util.JsonResponse(w, http.StatusUnauthorized, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
