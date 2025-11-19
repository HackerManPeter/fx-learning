package config

type cache struct {
	Host   string
	Port   string
	Prefix string
}

func newCache() cache {
	return cache{
		Host:   getEnv("REDIS_HOST", "127.0.0.1"),
		Port:   getEnv("REDIS_PORT", "6379"),
		Prefix: getEnv("REDIS_PREFIX", "fx_tutorial"),
	}
}
