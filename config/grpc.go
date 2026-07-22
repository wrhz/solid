package config

type GrpcConfig struct {
	useGrpc bool

	host string
	port int

	serverTlsCertFile       string
	serverTlsKeyFile        string

	clientTlsCertFile       string
	clientTlsKeyFile        string
	clientCaCertFile        string
}

func NewGrpcConfig() *GrpcConfig {
	return &GrpcConfig{
		useGrpc: false,
		host:    "localhost",
		port:    50051,
	}
}

func (g *GrpcConfig) UseGrpc(useGrpc bool) {
	g.useGrpc = useGrpc
}

func (g *GrpcConfig) GetUseGrpc() bool {
	return g.useGrpc
}

func (g *GrpcConfig) SetHost(host string) {
	g.host = host
}

func (g *GrpcConfig) GetHost() string {
	return g.host
}

func (g *GrpcConfig) SetPort(port int) {
	g.port = port
}

func (g *GrpcConfig) GetPort() int {
	return g.port
}

func (g *GrpcConfig) SetServerTlsCertFile(certFile string) {
	g.serverTlsCertFile = certFile
}

func (g *GrpcConfig) GetServerTlsCertFile() string {
	return g.serverTlsCertFile
}

func (g *GrpcConfig) SetServerTlsKeyFile(keyFile string) {
	g.serverTlsKeyFile = keyFile
}

func (g *GrpcConfig) GetServerTlsKeyFile() string {
	return g.serverTlsKeyFile
}

func (g *GrpcConfig) SetClientTlsCertFile(certFile string) {
	g.clientTlsCertFile = certFile
}

func (g *GrpcConfig) GetClientTlsCertFile() string {
	return g.clientTlsCertFile
}

func (g *GrpcConfig) SetClientTlsKeyFile(keyFile string) {
	g.clientTlsKeyFile = keyFile
}

func (g *GrpcConfig) GetClientTlsKeyFile() string {
	return g.clientTlsKeyFile
}

func (g *GrpcConfig) SetClientCaCertFile(CaCertFile string) {
	g.clientCaCertFile = CaCertFile
}

func (g *GrpcConfig) GetClientCaCertFile() string {
	return g.clientCaCertFile
}