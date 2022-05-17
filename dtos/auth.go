package dtos

type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
	Role     string `form:"role" json:"role" binding:"required,oneof=admin user"`
}

type LoginResponse struct {
}

type RegisterForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
	Name     string `form:"name" json:"name" binding:"required,min=3,max=50"`
	Role     string `form:"role" json:"role" binding:"required,min=3,max=50"`
}

type RefreshTokenForm struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}
