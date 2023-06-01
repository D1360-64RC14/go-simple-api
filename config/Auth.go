package config

type Auth struct {
	Base64TokenSeed string `json:"base64TokenSeed"`
	BCryptCost      uint   `json:"bcryptCost"`
}
