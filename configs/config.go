package configs

// TODO: map[string]string => struct olarak config üzerinden oku

type Config struct {
	Server struct {
		Port string
		Host string
	}
	Database struct {
		Connection     string
		DatabaseName   string
		CollectionName string
	}
}

var Configuration = map[string]Config{
	"test": {
		Server:   {Port: "8080", Host: "localhost"},
		Database: {CollectionName:"mongodb://localhost:27017",DatabaseName: "booksDB", CollectionName: "books"}
		},
	//"qa":   {},
	//"prod": {},
}

func GetConfig(env string) map[string]string {
	if conf, ok := configs[env]; ok {
		return conf
	}
	return configs["test"]

}
