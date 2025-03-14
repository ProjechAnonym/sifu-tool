package models

import "github.com/golang-jwt/jwt/v5"

type Jwt struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Secret   string `yaml:"secret"`
}