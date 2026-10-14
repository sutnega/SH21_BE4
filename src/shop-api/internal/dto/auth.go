package dto

// RegisterRequest — тело запроса на регистрацию.
type RegisterRequest struct {
	// Email пользователя
	Email string `json:"email"    validate:"required,email"           example:"test@example.com"`
	// Имя
	Name string `json:"name"     validate:"required,min=1,max=100"   example:"Ivan"`
	// Фамилия
	Surname string `json:"surname"  validate:"required,min=1,max=100"   example:"Ivanov"`
	// Телефон
	Phone string `json:"phone"    validate:"required,min=5,max=20"    example:"+79001234567"`
	// Пароль (минимум 6 символов)
	Password string `json:"password" validate:"required,min=6"          example:"secret123"`
}

// LoginRequest — тело запроса на логин.
type LoginRequest struct {
	// Email пользователя
	Email string `json:"email"    validate:"required,email"   example:"test@example.com"`
	// Пароль
	Password string `json:"password" validate:"required"        example:"secret123"`
}

// ResetRequest — тело запроса на восстановление пароля.
type ResetRequest struct {
	// Email пользователя
	Email string `json:"email" validate:"required,email" example:"test@example.com"`
}

// TokenResponse — тело ответа с токеном.
type TokenResponse struct {
	// JWT-токен
	Token string `json:"token"   example:"eyJhbGciOiJIUzI1NiIs..."`
	// UUID пользователя
	UserID string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// StatusResponse — универсальный ответ для методов без тела (reset).
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}
