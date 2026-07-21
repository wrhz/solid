package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/wrhz/solid/database"
)

func WebSocketMiddleware(handler http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "requestId", id)
		w.Header().Set("X-Request-ID", id)

		fmt.Printf("[%s] WebSocket %s\n", time.Now().Format("2006-01-02 15:04:05"), r.URL.Path)

		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}