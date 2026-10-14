package dto

// ErrorResponse — универсальное тело ошибки.
type ErrorResponse struct {
	// Текст ошибки
	Error string `json:"error" example:"invalid request body"`
}

// AddressDTO — адрес клиента или поставщика.
type AddressDTO struct {
	// Страна
	Country string `json:"country" validate:"required,min=1,max=100" example:"Russia"`
	// Город
	City string `json:"city"    validate:"required,min=1,max=100" example:"Moscow"`
	// Улица и номер дома
	Street string `json:"street"  validate:"required,min=1,max=200" example:"Tverskaya 1"`
}

// ---------- Client ----------

// CreateClientRequest — тело запроса на создание клиента.
type CreateClientRequest struct {
	// Имя клиента
	ClientName string `json:"client_name"    validate:"required,min=1,max=100"      example:"Ivan"`
	// Фамилия клиента
	ClientSurname string `json:"client_surname" validate:"required,min=1,max=100"      example:"Ivanov"`
	// Дата рождения в формате YYYY-MM-DD
	Birthday string `json:"birthday"       validate:"required,datetime=2006-01-02" example:"1990-05-15"`
	// Пол: male, female или other
	Gender string `json:"gender"         validate:"required,oneof=male female other" example:"male"`
	// Адрес клиента
	Address AddressDTO `json:"address"        validate:"required"`
}

// ClientResponse — тело ответа с данными клиента.
type ClientResponse struct {
	// Уникальный идентификатор клиента
	ID string `json:"id"                example:"550e8400-e29b-41d4-a716-446655440000"`
	// Имя клиента
	ClientName string `json:"client_name"       example:"Ivan"`
	// Фамилия клиента
	ClientSurname string `json:"client_surname"    example:"Ivanov"`
	// Дата рождения в формате YYYY-MM-DD
	Birthday string `json:"birthday"          example:"1990-05-15"`
	// Пол
	Gender string `json:"gender"            example:"male"`
	// Дата регистрации в формате RFC3339
	RegistrationDate string `json:"registration_date" example:"2024-01-01T12:00:00Z"`
	// Адрес клиента
	Address AddressDTO `json:"address"`
}

// UpdateAddressRequest — тело запроса на изменение адреса.
type UpdateAddressRequest struct {
	// Страна
	Country string `json:"country" validate:"required,min=1,max=100" example:"Russia"`
	// Город
	City string `json:"city"    validate:"required,min=1,max=100" example:"Moscow"`
	// Улица и номер дома
	Street string `json:"street"  validate:"required,min=1,max=200" example:"Tverskaya 1"`
}

// ---------- Product ----------

// CreateProductRequest — тело запроса на создание товара.
type CreateProductRequest struct {
	// Название товара
	Name string `json:"name"            validate:"required,min=1,max=200" example:"Bosch KGN39"`
	// Название категории (если её нет, будет создана)
	Category string `json:"category"        validate:"required,min=1,max=100" example:"Fridges"`
	// Цена
	Price float64 `json:"price"           validate:"required,gte=0"         example:"49999.99"`
	// Количество на складе
	AvailableStock int `json:"available_stock" validate:"gte=0"                  example:"10"`
	// UUID поставщика
	SupplierID string `json:"supplier_id"     validate:"required,uuid"          example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ProductResponse — тело ответа с данными товара.
type ProductResponse struct {
	// Уникальный идентификатор товара
	ID string `json:"id"               example:"550e8400-e29b-41d4-a716-446655440000"`
	// Название товара
	Name string `json:"name"             example:"Bosch KGN39"`
	// Категория
	Category string `json:"category"         example:"Fridges"`
	// Цена
	Price float64 `json:"price"            example:"49999.99"`
	// Количество на складе
	AvailableStock int `json:"available_stock"  example:"10"`
	// Дата последнего изменения остатка в формате RFC3339
	LastUpdateDate string `json:"last_update_date" example:"2024-01-01T12:00:00Z"`
	// UUID поставщика
	SupplierID string `json:"supplier_id"      example:"550e8400-e29b-41d4-a716-446655440000"`
	// UUID изображения (может отсутствовать)
	ImageID *string `json:"image_id,omitempty"`
}

// DecreaseStockRequest — тело запроса на уменьшение остатка товара.
type DecreaseStockRequest struct {
	// На сколько уменьшить остаток
	Amount int `json:"amount" validate:"required,gt=0" example:"1"`
}

// ---------- Supplier ----------

// CreateSupplierRequest — тело запроса на создание поставщика.
type CreateSupplierRequest struct {
	// Название поставщика
	Name string `json:"name"         validate:"required,min=1,max=200" example:"Bosch"`
	// Адрес поставщика
	Address AddressDTO `json:"address"      validate:"required"`
	// Телефон
	PhoneNumber string `json:"phone_number" validate:"required,min=5,max=20" example:"+7-800-555-35-35"`
}

// SupplierResponse — тело ответа с данными поставщика.
type SupplierResponse struct {
	// Уникальный идентификатор поставщика
	ID string `json:"id"           example:"550e8400-e29b-41d4-a716-446655440000"`
	// Название поставщика
	Name string `json:"name"         example:"Bosch"`
	// Адрес поставщика
	Address AddressDTO `json:"address"`
	// Телефон
	PhoneNumber string `json:"phone_number" example:"+7-800-555-35-35"`
}

// ---------- Image ----------

// ImageResponse — тело ответа после загрузки изображения.
type ImageResponse struct {
	// UUID изображения
	ID string `json:"id"         example:"550e8400-e29b-41d4-a716-446655440000"`
	// UUID товара, к которому привязано изображение
	ProductID string `json:"product_id" example:"660e8400-e29b-41d4-a716-446655440000"`
}
