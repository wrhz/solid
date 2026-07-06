package method

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wrhz/solid/database"
	"github.com/wrhz/solid/server"
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

type SolidMainRoute interface {
	SolidRoute

	ServerStart()
	ServerEnd()
}

func routeFuncHandle(callFunc func(c *server.Context) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := &server.Context{Writer: w, Request: req}

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

func NewRoute() *RouteStruct {
	return &RouteStruct{perfix: ""}
}

func (r *RouteStruct) Any(path string, callFunc func(c *server.Context) error) {
	r.Get(path, callFunc)
	r.Post(path, callFunc)
	r.Patch(path, callFunc)
	r.Delete(path, callFunc)
	r.Put(path, callFunc)
	r.Options(path, callFunc)
	r.Head(path, callFunc)
}

func (r *RouteStruct) Group(prefix string, callStruct SolidRoute) {
	route := &RouteStruct{perfix: r.perfix + prefix, middlewares: r.middlewares}
	callStruct.Init(route)
	callStruct.RegisterMiddleware(route)
	callStruct.RegisterRoute(route)
}

func (r *RouteStruct) Use(middleware func(c *server.Context, next http.HandlerFunc)) {
	r.middlewares = append(r.middlewares, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware(&server.Context{Writer: w, Request: r}, next.ServeHTTP)
		})
	})
}
