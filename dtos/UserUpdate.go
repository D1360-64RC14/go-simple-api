package dtos

type UserUpdate struct {
	UserName string `json:"username" binding:"min=3,max=50,alphanum"`
}
