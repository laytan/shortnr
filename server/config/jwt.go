package config

import "time"

// TokenDuration specifies how long a JWT token is valid
const TokenDuration = time.Minute * 15

// RefreshTokenDuration specifies how long a refresh token is valid (the time before a new login is required)
const RefreshTokenDuration = time.Hour * 48
