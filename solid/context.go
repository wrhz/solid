package solid

import (
	"net/http"

	"gorm.io/gorm"
)

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request

	gormDatabase *gorm.DB
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
