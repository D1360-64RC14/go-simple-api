package config

type Database struct {
	Address      string `yaml:"address"`
	DBName       string `yaml:"dbName"`
	Username     string `yaml:"username"`
	RootPassword string `yaml:"rootPassword"`
}
