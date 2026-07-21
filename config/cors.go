package config

type CorsConfig struct {
    useCors bool

	allowOrigin []string
	allowCredentials bool
	allowMethods []string
	allowHeaders []string
	exposeHeaders []string
	maxAge uint
}

func NewCorsConfig() *CorsConfig {
	return &CorsConfig{
        useCors: false,
        allowOrigin: []string{"*"},
        allowCredentials: false,
        allowMethods: []string{"GET"},
        allowHeaders: []string{},
        exposeHeaders: []string{},
        maxAge: 3600,
    }
}

func (c *CorsConfig) GetUseCors() bool {
    return c.useCors
}

func (c *CorsConfig) UseCors(useCors bool) {
    c.useCors = useCors
}

func (c *CorsConfig) GetAllowOrigin() []string {
    if c.allowOrigin == nil {
        c.allowOrigin = []string{}
    }

    return c.allowOrigin
}

func (c *CorsConfig) SetAllowOrigin(origins []string) {
    c.allowOrigin = origins
}

func (c *CorsConfig) GetAllowCredentials() bool {
    return c.allowCredentials
}

func (c *CorsConfig) SetAllowCredentials(allow bool) {
    c.allowCredentials = allow
}

func (c *CorsConfig) GetAllowMethods() []string {
    if c.allowMethods == nil {
        c.allowMethods = []string{}
    }

    return c.allowMethods
}

func (c *CorsConfig) SetAllowMethods(methods []string) {
    c.allowMethods = methods
}

func (c *CorsConfig) GetAllowHeaders() []string {
    if c.allowHeaders == nil {
        c.allowHeaders = []string{}
    }
    
    return c.allowHeaders
}

func (c *CorsConfig) SetAllowHeaders(headers []string) {
    c.allowHeaders = headers
}

func (c *CorsConfig) GetExposeHeaders() []string {
    if c.exposeHeaders == nil {
        c.exposeHeaders = []string{}
    }

    return c.exposeHeaders
}

func (c *CorsConfig) SetExposeHeaders(headers []string) {
    c.exposeHeaders = headers
}

func (c *CorsConfig) GetMaxAge() uint {
    return c.maxAge
}

func (c *CorsConfig) SetMaxAge(age uint) {
    c.maxAge = age
}