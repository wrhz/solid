package binding

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
)

type JSONBinding struct{}

func (JSONBinding) Name() string {
	return "json"
}

func (j JSONBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
        return fmt.Errorf("Content-Type must be application/json, got %s", ct)
    }

	defer req.Body.Close()

    body, err := io.ReadAll(req.Body)

    if err != nil {
		return err
	}

	return j.BindBody(body, s)
}

func (j JSONBinding) BindBody(body []byte, s any) error {
	return json.Unmarshal(body, s)
}