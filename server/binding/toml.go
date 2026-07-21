package binding

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
)

type TOMLBinding struct{}

func (TOMLBinding) Name() string {
	return "toml"
}

func (t TOMLBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/toml") {
        return fmt.Errorf("Content-Type must be application/xml, got %s", ct)
    }

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	return t.BindBody(body, s)
}

func (t TOMLBinding) BindBody(body []byte, s any) error {
	return toml.Unmarshal(body, s)
}