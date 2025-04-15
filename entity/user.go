package entity

type User struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	// password always keep hashed password
	Password string `json:"-"`
	Role     Role   `json:"role"`
}
