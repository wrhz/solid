package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wrhz/solid/database"

	solidManager "github.com/wrhz/solid/manager"
)

func RouteMiddleware(handler http.Handler)func (w http.ResponseWriter, r *http.Request)  {
	settings := solidManager.GetSettingsConfig()

	return func (w http.ResponseWriter, r *http.Request) {
		defer func() {
			id, ok := r.Context().Value("requestID").(string)

			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				httpError(w)
				debug.PrintStack()

				if !ok {
					return
				}

				if orm, ok := database.GormDatabasesManager.Get(id); ok {
					orm.Rollback()
				}
			} else {
				if !ok {
					return
				}

				if orm, ok := database.GormDatabasesManager.Get(id); ok {
					orm.Commit()
				}
			}
		}()

		corsConfig := solidManager.GetCorsConfig()

		if corsConfig.GetUseCors() && !corsMiddleware(w, r) {
			return
		}

		maxBytesMemory := settings.GetMaxBytesMemory()

		r.Body = http.MaxBytesReader(w, r.Body, maxBytesMemory)

		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			multipartFormMaxMemory := settings.GetMultipartFormMaxMemory()

			err := r.ParseMultipartForm(multipartFormMaxMemory)

			if err != nil {
				fmt.Println("Error parsing multipart form:", err)
				http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
				return
			}
		}

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "requestID", id)
		w.Header().Set("X-Request-ID", id)

		staticMaxAge := settings.GetStaticMaxAge()

		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", staticMaxAge))

		timeStart := time.Now()

		handler.ServeHTTP(w, r.WithContext(ctx))

		fmt.Printf("[%s] %s %s ... %v\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, time.Since(timeStart))
	}
}