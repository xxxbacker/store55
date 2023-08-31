package internal

type Config struct {
	Db_string string `yaml:"db_string"`
	Address   string `yaml:"address"`
}

func NewConfig() *Config {
	return &Config{
		Db_string: "",
		Address:   "",
	}
}
