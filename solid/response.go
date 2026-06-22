package solid

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func StringResponse(c *Context, s string, status int) error {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", s)

	return nil
}

func BytesResponse(c *Context, data []byte, status int) error {
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(status)

	c.Writer.Write(data)

	return nil
}

func JsonResponse(c *Context, data any, status int) error {
	var jsonData = gjson.New(data).MustToJsonString()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", jsonData)

	return nil
}

func HtmlResponse(c *Context, html string, status int) error {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", html)

	return nil
}

func HtmlViewResponse(c *Context, file string, status int) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "view", file + ".html"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", html)

	return nil
}

func VueViewResponse(c *Context, group string, status int) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "vue", group, group + ".html"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", html)

	return nil
}

func ReactViewResponse(c *Context, group string, status int) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "react", group, group + ".html"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", html)

	return nil
}

func XmlResponse(c *Context, data any, status int) error {
	var xmlData, err = xml.Marshal(data)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to marshal xml: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/xml")
	c.Writer.WriteHeader(status)

	fmt.Fprintf(c.Writer, "%s", xmlData)

	return nil
}

func (c *Context) Redirect(url string, status int) error {
	http.Redirect(c.Writer, c.Request, url, status)

	return nil
}

func (c *Context) NoContent() error {
	c.Writer.WriteHeader(http.StatusNoContent)

	return nil
}

func (c *Context) File(filePath string) error {
	http.ServeFile(c.Writer, c.Request, filePath)

	return nil
}

func (c *Context) Download(filePath string, fileName string) error {
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	http.ServeFile(c.Writer, c.Request, filePath)

	return nil
}

func (c *Context) Stream(streamFunc func(w http.ResponseWriter)) error {
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	streamFunc(c.Writer)

	return nil
}

func (c *Context) Error(status int, err error) error {
	c.Writer.WriteHeader(status)
	fmt.Fprintf(c.Writer, "%s", err.Error())

	return nil
}

func (c *Context) JSONError(status int, err error) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	
	JsonResponse(c, map[string]error{ "error": err }, status)

	return nil
}