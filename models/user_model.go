package models

type User struct {
	// gorm.Model
	UserID        int    `json:"user_id" gorm:"primaryKey"`
	UserLogin     string `json:"user_login"`
	UserPassword  string `json:"user_password"`
	UserCitizenId string `json:"user_citizen_id"`
	UserRole      int    `json:"user_role"`
	// Add more fields as needed
}
