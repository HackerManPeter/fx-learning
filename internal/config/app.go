package config

type Environment string

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentDevelopment Environment = "development"
)

type app struct {
	Port         string
	Environment  Environment
	JWTSecretKey string
}

func newApp() app {

	return app{
		Port:         getEnv("APP_PORT", "3000"),
		Environment:  Environment(getEnv("ENVIRONMENT", "development")),
		JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),
	}
}
