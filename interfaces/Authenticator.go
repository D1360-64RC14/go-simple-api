package interfaces

type Authenticator interface {
	GenerateToken(id int, email string) (string, error)
	IsTokenValid(inputToken string) (bool, error)
}
