package config

type bible struct {
	URL string
}

func newBible() bible {
	return bible{
		URL: getEnv("BIBLE_BASE_URL", "https://bible-api.com"),
	}
}
