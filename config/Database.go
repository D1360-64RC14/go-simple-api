package config

type Database struct {
	Url          string `yaml:"url"`
	DBName       string `yaml:"dbName"`
	RootPassword string `yaml:"rootPassword"`
}
