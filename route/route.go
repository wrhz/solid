package method

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/wrhz/solid/database"
	"github.com/wrhz/solid/server"
)

var routeNames = map[string]string{}

type RouteFunc func(c *server.Context) error
type MiddlewareFunc func(c *server.Context)

type RouteStruct struct {
	perfix      string
	middlewares []func(http.Handler) http.Handler
}

type MiddlewareStruct struct {
	middlewares *[]func(http.Handler) http.Handler
}

type SolidRoute interface {
	Init() error

	RegisterRoute(*RouteStruct)
	RegisterMiddleware(*MiddlewareStruct)
}

type SolidMainRoute interface {
	SolidRoute

	ServerStart() error
	ServerEnd() error
}

type Route struct {
	perfix string
}

func (r *Route) Name(name string) {
	routeNames[name] = r.perfix
}

func routeFuncHandle(callFunc RouteFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := server.NewContext(req, w)

		routeRequestStart(ctx)
		
		routeRequestEnd(ctx, callFunc(ctx))
	})
}

func routeRequestStart(ctx *server.Context) {
	id := ctx.RequestID()

	if id != "" {
		if database.IsStartGorm() {
			database.GormDatabasesManager.Set(id)

			db, ok := database.GormDatabasesManager.Get(id)

			if ok {
				ctx.SetGormDatabase(db)
			}
		}

		if database.IsStartXorm() {
			err := database.XormSessionsManager.Set(id)

			if err != nil {
				fmt.Println(err)
				return
			}

			session, ok := database.XormSessionsManager.Get(id)

			if ok {
				ctx.SetXormSession(session)
			}
		}
	}
}

func routeRequestEnd(ctx *server.Context, err error) {
	id := ctx.RequestID()

	if id == "" { return }

	if database.IsStartGorm() {
		tx, ok := database.GormDatabasesManager.Get(id)

		if !ok { return }

		if err == nil {
			if err = tx.Commit().Error; err != nil {
				log.Fatal("Commit error:", err)
			}
		} else {
			if err = tx.Rollback().Error; err != nil {
				log.Fatal("‌Rollback error:", err)
			}
		}

		database.GormDatabasesManager.Delete(id)
	}

	if database.IsStartXorm() {
		session, ok := database.XormSessionsManager.Get(id)

		if !ok { return }

		if err == nil {
			if err = session.Commit(); err != nil {
				log.Fatal("Commit error:", err)
			}
		} else {
			if err = session.Rollback(); err != nil {
				log.Fatal("‌Rollback error:", err)
			}
		}

		database.XormSessionsManager.Delete(id)
	}
}

func (r *RouteStruct) GetMiddlewares() *[]func(http.Handler) http.Handler {
	return &r.middlewares
}

func NewRoute() *RouteStruct {
	return &RouteStruct{ perfix: "", middlewares: []func(http.Handler) http.Handler{} }
}

func NewMiddleware(middlewares *[]func(http.Handler) http.Handler) *MiddlewareStruct {
	return &MiddlewareStruct{ middlewares: middlewares }
}

func (r *RouteStruct) Any(path string, callFunc RouteFunc) *Route {
	r.Get(path, callFunc)
	r.Post(path, callFunc)
	r.Patch(path, callFunc)
	r.Delete(path, callFunc)
	r.Put(path, callFunc)
	r.Options(path, callFunc)
	return r.Head(path, callFunc)
}

func (r *RouteStruct) Group(prefix string, callStruct SolidRoute) error {
	route := &RouteStruct{ perfix: r.perfix + prefix, middlewares: r.middlewares }
	middleware := &MiddlewareStruct{ middlewares: &route.middlewares }

	if err := callStruct.Init(); err != nil {
		return err
	}
	
	callStruct.RegisterMiddleware(middleware)
	callStruct.RegisterRoute(route)

	return nil
}

func (r *MiddlewareStruct) Use(middleware MiddlewareFunc) {
	*r.middlewares = append(*r.middlewares, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context := server.NewContext(r, w)

			context.SetNext(next.ServeHTTP)

			middleware(context)
		})
	})
}

func Reverse(name string, args ...any) string {
	perfix, ok := routeNames[name]

	if !ok {
		return ""
	}

	re := regexp.MustCompile(`\{([^}]*)\}`)
	i := 0

	return re.ReplaceAllStringFunc(perfix, func(match string) string {
		if i < len(args) {
			val := args[i]
			i++
			return fmt.Sprint(val)
		}
		return ""
	})
}
