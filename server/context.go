package server

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"

	config "github.com/wrhz/solid/manager"
	"gorm.io/gorm"
	"xorm.io/xorm"
)

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request

	Errors []error

	gormDatabase *gorm.DB
	xormSession *xorm.Session

	jsonpCallback string

	data map[string]any

	next http.HandlerFunc

	body []byte
	getBody bool
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

func (c *Context) SetJSONPCallBack(callback string) {
	c.jsonpCallback = callback
}

func (c *Context) GetJSONPCallback() string {
	return c.jsonpCallback
}

func (c *Context) Set(key string, value any) {
	c.data[key] = value
}

func (c *Context) Get(key string) any {
	return c.data[key]
}

func (c *Context) Copy() Context {
	return *c
}

func (c *Context) SetNext(next http.HandlerFunc) {
	c.next = next
}

func (c *Context) Next() {
	if c.next != nil {
		c.next(c.Writer, c.Request)
	}
}

func (c *Context) Abort() {
	c.next = nil
}

func (c *Context) Error(err error) {
	c.Errors = append(c.Errors, err)
}

func (c *Context) ClientIP() string {
	trustedProxies := config.GetSettingsConfig().GetTrustedProxies()
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        ip = c.Request.RemoteAddr
    }
	if trustedProxies != nil && slices.Contains(trustedProxies, ip) {
		xff := c.Request.Header.Get("X-Forwarded-For")

		if xff == "" {
			return ip
		}

		parts := strings.Split(xff, ",")

		for i := len(parts) - 1; i >= 0; i-- {
			part := parts[i]

			if part == "" {
				continue
			}

			if !slices.Contains(trustedProxies, part) {
				return part
			}
		}
	}

	return ip
}

func (c *Context) Header() map[string][]string {
	return c.Request.Header
}

func (c *Context) Body() ([]byte, error) {
	if c.getBody {
		return c.body, nil
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		return nil, err
	}

	c.getBody = true
	c.body = body

	return body, nil
}

func (c *Context) SaveUploadedFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request: r,
		Writer: w,
		jsonpCallback: "callback",
	}
}

func CreateTestContext(w http.ResponseWriter) *Context {
	return &Context{
		Writer: w,
		jsonpCallback: "callback",
	}
}
