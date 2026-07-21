package config

type ICorsConfig interface {
    GetUseCors() bool
    GetAllowCredentials() bool
    GetAllowHeaders() []string
    GetAllowMethods() []string
    GetAllowOrigin() []string
    GetExposeHeaders() []string
    GetMaxAge() uint
}