package dtos

type IdentifiedUserWithHash struct {
	IdentifiedUser
	Hash string `json:"hash"`
}
