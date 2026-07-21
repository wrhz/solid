package binding

import "net/http"

type Binding interface {
	Name() string
	Bind(req *http.Request, obj any) error
}

type BindingBody interface {
    Binding
    BindBody([]byte, any) error
}

var Bindings = map[string]Binding{
	"application/x-www-form-urlencoded": FormBinding{},
	"multipart/form-data": FormBinding{},
	"application/json": JSONBinding{},
	"application/x-msgpack": MsgPackBinding{},
	"application/Protobuf": ProtobufBinding{},
	"application/toml": TOMLBinding{},
	"application/xml": XMLBinding{},
	"application/yaml": YAMLBinding{},
}

var BindingBodys = map[string]BindingBody{
	"application/json": JSONBinding{},
	"application/x-msgpack": MsgPackBinding{},
	"application/Protobuf": ProtobufBinding{},
	"application/toml": TOMLBinding{},
	"application/xml": XMLBinding{},
	"application/yaml": YAMLBinding{},
}