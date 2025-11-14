package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/source-build/mcp-pm/internal/config"
)

type JWTClaims struct {
	// 用户ID
	UserID string `json:"user_id"`
	// 项目ID
	ProjectID string `json:"project_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userId string, projectId string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:    userId,
		ProjectID: projectId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // Token expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // Token issuance time
			NotBefore: jwt.NewNumericDate(time.Now()),               // Token validity start time
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.ServerJwtSecret))
}

// VerifyJWT 校验JWT token
func VerifyJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Verify the signing method is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.ServerJwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}

	// 校验token claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("%w: invalid token claims", auth.ErrInvalidToken)
}
