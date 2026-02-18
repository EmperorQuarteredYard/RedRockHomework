package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	AccessSecret  string        `json:"access_secret"`
	RefreshSecret string        `json:"refresh_secret"`
	Issuer        string        `json:"issuer"`
	AccessTTL     time.Duration `json:"access_ttl"`
	RefreshTTL    time.Duration `json:"refresh_ttl"`
}

var (
	config            *JWTConfig
	configInitialized bool = false
)

type CustomClaims struct {
	UserID     uint64 `json:"user_id"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Type       string `json:"type"` // "access"/"refresh"
	jwt.RegisteredClaims
}

type AuthUser struct {
	UserID     uint64 `json:"user_id"`
	Role       string `json:"role"`
	Department string `json:"department"`
	//Nickname   string `json:"nickname"` 不加了要改好多石山
}

var once sync.Once

func getConfig() (err error) {
	if configInitialized {
		return nil
	}

	newConfig := JWTConfig{}
	var configFile *os.File
	configFile, err = os.Open("backend/config/jwt.json")

	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&newConfig)
	if err != nil {
		return
	}
	config = &newConfig

	configInitialized = true
	return err
}

func GenerateRefreshToken(userID uint64, department, role string) (token string, err error) {
	err = getConfig()
	if err != nil {
		return "", err
	}
	now := time.Now()
	refreshClaims := CustomClaims{
		UserID:     userID,
		Department: department,
		Role:       role,
		Type:       "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(now.Add(config.RefreshTTL)),
			NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshTok := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	token, err = refreshTok.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateAccessToken(userID uint64, department, role string) (token string, err error) {
	err = getConfig()
	if err != nil {
		return "", err
	}
	now := time.Now()
	accessClaims := CustomClaims{
		UserID:     userID,
		Department: department,
		Role:       role,
		Type:       "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(now.Add(config.AccessTTL)),
			NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessTok := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	token, err = accessTok.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateToken(userID uint64, department, role string) (accessToken string, refreshToken string, err error) {
	refreshToken, err = GenerateRefreshToken(userID, department, role)
	if err != nil {
		return "", "", err
	}

	accessToken, err = GenerateAccessToken(userID, department, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func VerifyAccessToken(accessToken string) (claims *CustomClaims, err error) {
	if !configInitialized {
		err = getConfig()
		if err != nil {
			return nil, err
		}
	}
	raw := stripBearer(accessToken)
	token, err := jwt.ParseWithClaims(raw, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.AccessSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid signing token")
	}
	if claims.Type != "access" {
		return nil, fmt.Errorf("token type mismatch: not access token")
	}
	return claims, nil
}

func VerifyRefreshToken(refreshToken string) (claims *CustomClaims, err error) {
	if !configInitialized {
		err = getConfig()
		if err != nil {
			return nil, err
		}
	}
	raw := stripBearer(refreshToken)
	token, err := jwt.ParseWithClaims(raw, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.RefreshSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("token type mismatch: not refresh token")
	}
	return claims, nil
}

func stripBearer(s string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(s)), "bearer ") {
		return strings.TrimSpace(s[len("Bearer "):])
	}
	return strings.TrimSpace(s)
}

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header required",
				"code":    StatusMissedToken,
				"data":    nil,
			})
			return
		}

		claims, err := VerifyAccessToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
				"code":    StatusInvalidToken,
				"data":    nil})
			return
		}

		// 将用户信息存入上下文
		c.Set("user", &AuthUser{
			UserID:     claims.UserID,
			Role:       claims.Role,
			Department: claims.Department,
		})
		c.Next()
	}
}

// RoleMiddleware 角色检查中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    StatusUserNotAuthenticated,
				"message": "User not authenticated",
				"data":    nil,
			})
			return
		}

		user, ok := userInterface.(*AuthUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    StatusInvalidUserData,
				"message": "Invalid user data",
				"data":    nil,
			})
			return
		}

		// 检查用户角色是否在允许的角色列表中
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if user.Role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Insufficient permissions",
				"code":    StatusInsufficientPermissions,
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(c *gin.Context) (*AuthUser, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	user, ok := userInterface.(*AuthUser)
	if !ok {
		return nil, fmt.Errorf("invalid user data in context")
	}

	return user, nil
}
