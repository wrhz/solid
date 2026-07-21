package binding

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackBinding struct{}

func (MsgPackBinding) Name() string {
	return "msgpack"
}

func (m MsgPackBinding) Bind(req *http.Request, s any) error {
	if ct := req.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-msgpack") {
        return fmt.Errorf("Content-Type must be application/xml, got %s", ct)
    }

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	return m.BindBody(body, s)
}

func (m MsgPackBinding) BindBody(body []byte, s any) error {
	return msgpack.Unmarshal(body, s)
}