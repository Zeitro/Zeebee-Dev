package llm

const envPrefix = "env-"

// EnvFilename builds filename for given env and name
func EnvFilename(env, name string) string {
	if env == "" {
		return name
	}
	return envPrefix + env + "-" + name
}
