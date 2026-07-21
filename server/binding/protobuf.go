package binding

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/protobuf/proto"
)

type ProtobufBinding struct{}

func (ProtobufBinding) Name() string {
	return "protobuf"
}

func (p ProtobufBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/protobuf") {
        return fmt.Errorf("Content-Type must be application/xml, got %s", ct)
    }

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	return p.BindBody(body, s)
}

func (p ProtobufBinding) BindBody(body []byte, s any) error {
	return proto.Unmarshal(body, s.(proto.Message))
}