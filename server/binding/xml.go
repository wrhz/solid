package binding

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type XMLBinding struct{}

func (XMLBinding) Name() string {
	return "xml"
}

func (x XMLBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/xml") {
        return fmt.Errorf("Content-Type must be application/xml, got %s", ct)
    }

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return fmt.Errorf("Failed to read: %w", err)
	}

	return x.BindBody(body, s)
}

func (x XMLBinding) BindBody(body []byte, s any) error {
	return xml.Unmarshal(body, s)
}