package config

import (
	"crypto/tls"
	"net/http"

	solidRoute "github.com/wrhz/solid/route"
)

type ServerConfigStruct struct {
	host       string
	port       int
	mainStruct solidRoute.SolidMainRoute

	tlsCertFile string
	tlsKeyFile  string
	tlsConfig   *tls.Config

	serverConfig *http.Server

	debug bool
}

func NewServerConfig() *ServerConfigStruct {
	return &ServerConfigStruct{
		host:       "localhost",
		port:       8000,
		mainStruct: nil,
		debug:      true,
	}
}

func (s *ServerConfigStruct) SetPort(port int) {
	s.port = port
}

func (s *ServerConfigStruct) SetHost(host string) {
	s.host = host
}

func (s *ServerConfigStruct) SetMainStruct(mainStruct solidRoute.SolidMainRoute) {
	if mainStruct == nil {
		panic("The mainStruct can't be nil")
	}

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

func (s *ServerConfigStruct) SetServerConfig(serverConfig *http.Server) {
	s.serverConfig = serverConfig
}

func (s *ServerConfigStruct) GetPort() int {
	return s.port
}

func (s *ServerConfigStruct) GetHost() string {
	return s.host
}

func (s *ServerConfigStruct) GetMainStruct() solidRoute.SolidMainRoute {
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

func (s *ServerConfigStruct) GetServerConfig() *http.Server {
	if s.serverConfig == nil {
		s.serverConfig = &http.Server{}
	}

	return s.serverConfig
}

func (s *ServerConfigStruct) GetDebug() bool {
	return s.debug
}

func (s *ServerConfigStruct) SetDebug(debug bool) {
	s.debug = debug
}
