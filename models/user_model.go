package models

type User struct {
	// gorm.Model
	UserID     int    `json:"user_id" gorm:"primaryKey"`
	UserLogin  string `json:"user_login"`
	UserPwd    string `json:"user_pwd"`
	UserIdCard string `json:"user_IDcard"`
	// Add more fields as needed
}
