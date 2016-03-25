package dawg

func sampleConfig() *Config {
	cfg := &Config{
		"www.github.com": &ServiceConfig{
			Template: URITemplate{nil, "https://github.com/{org}/{repo}"},
			Keyword:  "dawg gh",
			Substitutions: map[string]map[string]interface{}{
				"dawg": map[string]interface{}{
					"org":  "v-yarotsky",
					"repo": "dawg",
				},
			},
		},
	}
	return cfg
}
