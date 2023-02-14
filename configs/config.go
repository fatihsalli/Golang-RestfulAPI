package configs

// TODO: map[string]string => struct olarak config üzerinden oku

var config = map[string]map[string]string{
	"test": {"MONGOURI": "mongodb://localhost:27017"},
	"qa":   {"MONGOURI": "mongodb://localhost:27017"},
	"prod": {"MONGOURI": "mongodb://localhost:27017"},
}

func GetConfig(env string) map[string]string {
	if conf, ok := config[env]; ok {
		return conf
	}
	return config["test"]
}
