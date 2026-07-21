package binding

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLBinding struct{}

func (YAMLBinding) Name() string {
	return "yaml"
}

func (y YAMLBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/yaml") {
        return fmt.Errorf("Content-Type must be application/xml, got %s", ct)
    }

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	return y.BindBody(body, s)
}

func (y YAMLBinding) BindBody(body []byte, s any) error {
	return yaml.Unmarshal(body, s)
}