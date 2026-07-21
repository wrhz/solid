package server

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-json"
	"github.com/gogf/gf/v2/encoding/gjson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	solidManager "github.com/wrhz/solid/manager"
)

func (c *Context) String(s string, status int) error {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)

	_, err := fmt.Fprintf(c.Writer, "%s", s)

	return err
}

func (c *Context) Bytes(data []byte, status int) error {
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(status)

	_, err := c.Writer.Write(data)

	return err
}

func (c *Context) JSON(data any, status int) error {
	var jsonData = gjson.New(data).MustToJsonString()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	_, err := fmt.Fprintf(c.Writer, "%s", jsonData)

	return err
}

func (c *Context) HTML(html string, status int) error {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	_, err := fmt.Fprintf(c.Writer, "%s", html)

	return err
}

func (c *Context) View(file string, status int, args any) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "views", file))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	templateConfig := solidManager.GetTemplateConfig().GetTemplateRender()

	if templateConfig != nil {
		return templateConfig.Render(file, c.Writer, args)
	}

	_, err = c.Writer.Write(html)

	return err
}

func (c *Context) VueView(group string, status int) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "src", group, group + ".html"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprintf(c.Writer, "%s", html)

	return err
}

func (c *Context) ReactView(group string, status int) error {
	var html, err = os.ReadFile(filepath.Join(".", "resource", "src", group, group + ".html"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to read html file: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprintf(c.Writer, "%s", html)

	return err
}

func (c *Context) XML(data any, status int) error {
	var xmlData, err = xml.Marshal(data)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Failed to marshal xml: %s", err)
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/xml")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprintf(c.Writer, "%s", xmlData)

	return err
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

	_, err = fmt.Fprintf(c.Writer, "%s", err.Error())

	return err
}

func (c *Context) JSONError(status int, err error) error {
	c.Writer.WriteHeader(status)

	return c.String(fmt.Sprintf("{ \"error\": \"%s\" }", err), status)
}

func (c *Context) PureJSON(s any, status int) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	buf := new(bytes.Buffer)

    enc := json.NewEncoder(buf)

    enc.SetEscapeHTML(false)

    err := enc.Encode(s)

    if err != nil {
        return err
    }

    _, err = c.Writer.Write(bytes.TrimRight(buf.Bytes(), "\n"))

	return err
}

func (c *Context) AsciiJSON(s any, status int) error {
	var jsonData = gjson.New(s).MustToJsonString()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	quoted := strconv.QuoteToASCII(jsonData)

	_, err := fmt.Fprintf(c.Writer, "%s", quoted[1 : len(quoted)-1])

	return err
}

func (c *Context) SecureJSON(s any, status int) error {
	var jsonData = gjson.New(s).MustToJsonString()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	_, err := fmt.Fprintf(c.Writer, "%s", "while(1);" + jsonData)

	return err
}

func (c *Context) JSONP(s any, status int) error {
	var jsonData = gjson.New(s).MustToJsonString()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	callback := c.Request.URL.Query().Get(c.GetJSONPCallback())

	_, err := fmt.Fprintf(c.Writer, "%s", callback + "(" + jsonData + ")")

	return err
}

func (c *Context) YAML(s any, status int) error {
	data, err := yaml.Marshal(s)

	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/yaml")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprint(c.Writer, data)

	return err
}

func (c *Context) ProtoBuf(s proto.Message, status int) error {
	data, err := proto.Marshal(s)

	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/protobuf")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprint(c.Writer, data)

	return err
}

func (c *Context) TOML(s any, status int) error {
	data, err := toml.Marshal(s)

	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/toml")
	c.Writer.WriteHeader(status)

	_, err = fmt.Fprint(c.Writer, data)

	return err
}