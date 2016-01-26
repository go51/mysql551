package mysql551

type Mysql struct {
	config *Config
}

type Config struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func New(config *Config) *Mysql {
	m := Mysql{
		config:config,
	}

	return &m
}
