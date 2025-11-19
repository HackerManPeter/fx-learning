package config

type database struct {
	Username string
	Password string
	Port     string
	Name     string
	Host     string
}

func newDatabase() database {
	return database{
		Username: getEnv("DATABASE_USERNAME", "postgres"),
		Password: getEnv("DATABASE_PASSWORD", "postgres"),
		Port:     getEnv("DATABASE_PORT", "5432"),
		Name:     getEnv("DATABASE_NAME", "fx_tutorial"),
		Host:     getEnv("DATABASE_HOST", "127.0.0.1"),
	}
}
