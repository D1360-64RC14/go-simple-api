package models

type UserModel struct {
	UserName string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required,email,max=100"`
}
