package config

type Auth struct {
	Base64TokenSeed string `json:"base64TokenSeed"`
	BCryptCost      int    `json:"bcryptCost"`
}
