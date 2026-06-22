package solid

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RouteStruct struct {
	perfix      string
	middlewares []func(http.Handler) http.Handler
}

type SolidRoute interface {
	Init(*RouteStruct)

	RegisterRoute(*RouteStruct)
	RegisterMiddleware(*RouteStruct)
}

var getRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var postRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var patchRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var deleteRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var putRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var optionsRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var headRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}
var websocketRoutes = map[string]func(w http.ResponseWriter, r *http.Request) {}

func routeRequestStart(ctx *Context) {
	id := ctx.RequestID()

	if id != "" {
		GormDatabasesManager.Set(id)

		db, ok := GormDatabasesManager.Get(id)

		if ok {
			ctx.SetGormDatabase(db)
		}
	}
}

func routeRequestEnd(ctx *Context, err error) {
	id := ctx.RequestID()

	if id == "" { return }

	tx, ok := GormDatabasesManager.Get(id)

	if !ok { return }

	if err == nil {
		if err := tx.Commit().Error; err != nil {
			log.Fatal("Commit error:", err)
		}
	} else {
		if err := tx.Rollback().Error; err != nil {
			log.Fatal("‌Rollback error:", err)
		}
	}

	GormDatabasesManager.Delete(id)
}

func GetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return getRoutes
}

func PostRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return postRoutes
}

func PatchRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return patchRoutes
}

func DeleteRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return deleteRoutes
}

func PutRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return putRoutes
}

func OptionsRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return optionsRoutes
}

func HeadRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return headRoutes
}

func WebsocketRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	return websocketRoutes
}

func NewRoute() *RouteStruct {
	return &RouteStruct{perfix: ""}
}

func (r *RouteStruct) Get(path string, callFunc func(c *Context) error) {
	getRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}

		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Post(path string, callFunc func(c *Context) error) {
	postRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Patch(path string, callFunc func(c *Context) error) {
	patchRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Delete(path string, callFunc func(c *Context) error) {
	deleteRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Put(path string, callFunc func(c *Context) error) {
	putRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Options(path string, callFunc func(c *Context) error) {
	optionsRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Head(path string, callFunc func(c *Context) error) {
	headRoutes[r.perfix+path] = r.routeChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{Writer: w, Request: req}
		
		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})).ServeHTTP
}

func (r *RouteStruct) Any(path string, callFunc func(c *Context) error) {
	r.Get(path, callFunc)
	r.Post(path, callFunc)
	r.Patch(path, callFunc)
	r.Delete(path, callFunc)
	r.Put(path, callFunc)
	r.Options(path, callFunc)
	r.Head(path, callFunc)
}

func (r *RouteStruct) Websocket(path string, callFunc func(websocket *WebSocket)) {
	websocketConfig := GetWebSocketConfig()
	
	websocketRoutes[r.perfix+path] = r.websocketChain(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocketConfig.Upgrader.Upgrade(w, req, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}

		websocketStruct := &WebSocket{ Request: req, Conn: conn, isClosed: false }
		
		defer websocketStruct.Close()

		ticker := time.NewTicker(time.Duration(websocketConfig.GetPingInterval()) * time.Second)

		defer ticker.Stop()

		go func() {
			for range ticker.C {
				err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Duration(websocketConfig.GetPongWait()) * time.Second))
				if err != nil {
					fmt.Println("WebSocket ping error:", err)
					websocketStruct.Close()
					return
				}
			}
		}()

		callFunc(websocketStruct)
	})).ServeHTTP
}

func (r *RouteStruct) Group(prefix string, callStruct SolidRoute) {
	route := &RouteStruct{perfix: r.perfix + prefix, middlewares: r.middlewares}
	callStruct.Init(route)
	callStruct.RegisterMiddleware(route)
	callStruct.RegisterRoute(route)
}

func (r *RouteStruct) Use(middleware func(c *Context, next http.HandlerFunc)) {
	r.middlewares = append(r.middlewares, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware(&Context{Writer: w, Request: r}, next.ServeHTTP)
		})
	})
}

func (r *RouteStruct) routeChain(handler http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	settings := GetSettingsConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			id, ok := r.Context().Value("requestID").(string)

			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				if !ok { return }

				if orm, ok := GormDatabasesManager.Get(id); ok {
					orm.Rollback()
				}
			} else {
				if !ok { return }

				if orm, ok := GormDatabasesManager.Get(id); ok {
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
