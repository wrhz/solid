package method

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

func (r *RouteStruct) routeChain(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	settings := solidManager.GetSettingsConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			id, ok := r.Context().Value("requestID").(string)

			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

		maxBytesMemory, err := settings.GetMaxBytesMemory()
		if err != nil {
			fmt.Println("Error getting max bytes memory:", err)
			http.Error(w, "Error getting max bytes memory", http.StatusInternalServerError)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBytesMemory)

		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			multipartFormMaxMemory, err := settings.GetMultipartFormMaxMemory()
			if err != nil {
				fmt.Println("Error getting multipart form max memory:", err)
				http.Error(w, "Error getting multipart form max memory", http.StatusInternalServerError)
				return
			}

			err = r.ParseMultipartForm(multipartFormMaxMemory)

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

		staticMaxAge, err := settings.GetStaticMaxAge()
		if err != nil {
			fmt.Println("Error getting static max age:", err)
		}

		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", staticMaxAge))

		timeStart := time.Now()

		handler.ServeHTTP(w, r.WithContext(ctx))

		fmt.Printf("[%s] %s %s ... %v\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, time.Since(timeStart))
	})
}

func (r *RouteStruct) websocketChain(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "requestId", id)
		w.Header().Set("X-Request-ID", id)

		fmt.Printf("[%s] WebSocket %s\n", time.Now().Format("2006-01-02 15:04:05"), r.URL.Path)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}