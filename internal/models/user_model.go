package models

type UserResponse struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt uint8  `json:"created_at,omitempty"`
	UpdatedAt uint8  `json:"updated_at,omitempty"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
type SignUpRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GetTokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}
