package server

import (
	"net/http"

	"gorm.io/gorm"
	"xorm.io/xorm"
)

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request

	gormDatabase *gorm.DB
	xormSession *xorm.Session
}

func (c *Context) RequestID() string {
	if id, ok := c.Request.Context().Value("requestID").(string); ok {
		return id
	}

	return ""
}

func (c *Context) SetGormDatabase(db *gorm.DB) {
	c.gormDatabase = db
}

func  (c *Context) GetGormDatabase() *gorm.DB {
	return c.gormDatabase
}

func (c *Context) SetXormSession(session *xorm.Session) {
	c.xormSession = session
}

func (c *Context) GetXormSession() *xorm.Session {
	return c.xormSession
}
