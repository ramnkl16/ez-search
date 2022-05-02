package catalogboltdb

type Config struct {
	CoreDb    string `json:"coredb"`
	BuketName string `json:"buketName"`
}

var (
	config Config
)

func SetConfig(c Config) {
	config = c
}
