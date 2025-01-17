// internal/model/user.go
package model 

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-"`
    Name     string `json:"name"`
}