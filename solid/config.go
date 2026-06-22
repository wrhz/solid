package solid

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

var serverConfig = NewServerConfig()
var settings = NewSettingsConfig()
var websocketConfig = NewWebSocketConfig()
var databaseConfig = NewDatabaseConfig()

type ServerConfigStruct struct {
	host       string
	port       int
	mainStruct SolidRoute

	tlsCertFile string
	tlsKeyFile  string
	tlsConfig	*tls.Config
}

func NewServerConfig() *ServerConfigStruct {
	return &ServerConfigStruct{
		host:       "localhost",
		port:       8000,
		mainStruct: nil,
	}
}

func (s *ServerConfigStruct) SetPort(port int) {
	s.port = port
}

func (s *ServerConfigStruct) SetHost(host string) {
	s.host = host
}

func (s *ServerConfigStruct) SetMainStruct(mainStruct SolidRoute) {
	s.mainStruct = mainStruct
}

func (s *ServerConfigStruct) SetTLSCertFile(certFile string) {
	s.tlsCertFile = certFile
}

func (s *ServerConfigStruct) SetTLSKeyFile(keyFile string) {
	s.tlsKeyFile = keyFile
}

func (s *ServerConfigStruct) SetTLSConfig(tlsConfig *tls.Config) {
	s.tlsConfig = tlsConfig
}

func (s *ServerConfigStruct) GetPort() int {
	return s.port
}

func (s *ServerConfigStruct) GetHost() string {
	return s.host
}

func (s *ServerConfigStruct) GetMainStruct() SolidRoute {
	return s.mainStruct
}

func (s *ServerConfigStruct) GetTLSCertFile() string {
	return s.tlsCertFile
}

func (s *ServerConfigStruct) GetTLSKeyFile() string {
	return s.tlsKeyFile
}

func (s *ServerConfigStruct) GetTLSConfig() *tls.Config {
	return s.tlsConfig
}

func GetServerConfig() *ServerConfigStruct {
	return serverConfig
}

type SettingsConfigStruct struct {
	staticMaxAge			int

	maxBytesMemory         int64
	multipartFormMaxMemory int64

	sessionsPairs []byte
	sessionStore   sessions.Store
}

func NewSettingsConfig() *SettingsConfigStruct {
	return &SettingsConfigStruct{
		maxBytesMemory:         64 << 20,
		multipartFormMaxMemory: 32 << 20,
	}
}

func (s *SettingsConfigStruct) SetStaticMaxAge(staticMaxAge int) {
	s.staticMaxAge = staticMaxAge
}

func (s *SettingsConfigStruct) SetMaxBytesMemory(maxBytesMemory int64) {
	s.maxBytesMemory = maxBytesMemory
}

func (s *SettingsConfigStruct) SetMultipartFormMaxMemory(maxMemory int64) {
	s.multipartFormMaxMemory = maxMemory
}

func (s *SettingsConfigStruct) SetSessionsSecret(sessionsPairs ...string) {
	for _, pair := range sessionsPairs {
		s.sessionsPairs = append(s.sessionsPairs, []byte(pair)...)
	}
}

func (s *SettingsConfigStruct) SetSessionStore(sessionStore sessions.Store) {
	s.sessionStore = sessionStore
}

func (s *SettingsConfigStruct) GetStaticMaxAge() (int, error) {
	staticMaxAge := s.staticMaxAge
	if staticMaxAge == 0 {
		value, exists := os.LookupEnv("SOLID_STATIC_MAX_AGE")
		if exists {
			if val, err := strconv.Atoi(value); err == nil {
				staticMaxAge = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return staticMaxAge, nil
}

func (s *SettingsConfigStruct) GetMaxBytesMemory() (int64, error) {
	maxBytesMemory := s.maxBytesMemory
	if maxBytesMemory == 0 {
		value, exists := os.LookupEnv("SOLID_MAX_BYTES_MEMORY")
		if exists {
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				maxBytesMemory = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return maxBytesMemory, nil
}

func (s *SettingsConfigStruct) GetMultipartFormMaxMemory() (int64, error) {
	multipartFormMaxMemory := s.multipartFormMaxMemory
	if multipartFormMaxMemory == 0 {
		value, exists := os.LookupEnv("SOLID_MULTIPART_FORM_MAX_MEMORY")
		if exists {
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				multipartFormMaxMemory = val
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	}

	return multipartFormMaxMemory, nil
}

func (s *SettingsConfigStruct) GetSessionStore() sessions.Store {
	return s.sessionStore
}

func GetSettingsConfig() *SettingsConfigStruct {
	return settings
}

type WebSocketConfigStruct struct {
	Upgrader *Upgrader

	pingInterval int
	pongWait     int
}

func (w *WebSocketConfigStruct) SetPingInterval(pingInterval int) {
	w.pingInterval = pingInterval
}

func (w *WebSocketConfigStruct) SetPongWait(pongWait int) {
	w.pongWait = pongWait
}

func (w *WebSocketConfigStruct) GetPingInterval() int {
	return w.pingInterval
}

func (w *WebSocketConfigStruct) GetPongWait() int {
	return w.pongWait
}

func NewWebSocketConfig() *WebSocketConfigStruct {
	upgrader := &Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WebSocketConfigStruct{
		Upgrader: upgrader,
		pingInterval: 10,
		pongWait: 5,
	}
}

func GetWebSocketConfig() *WebSocketConfigStruct {
	return websocketConfig
}

type DatabaseConfigStruct struct {
	gormDialector gorm.Dialector
	gormOptions []gorm.Option
}

func (d *DatabaseConfigStruct) SetGormDialector(dialector gorm.Dialector) {
	d.gormDialector = dialector
}

func (d *DatabaseConfigStruct) SetGormOptions(options ...gorm.Option) {
	d.gormOptions = append(d.gormOptions, options...)
}

func (d *DatabaseConfigStruct) GetGormDialector() gorm.Dialector {
	return d.gormDialector
}

func (d *DatabaseConfigStruct) GetGormOptions() []gorm.Option {
	return d.gormOptions
}

func NewDatabaseConfig() *DatabaseConfigStruct {
	return &DatabaseConfigStruct{}
}

func GetDatabaseConfig() *DatabaseConfigStruct {
	return databaseConfig
}