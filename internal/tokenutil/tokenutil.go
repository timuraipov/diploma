package tokenutil

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/timuraipov/diploma/internal/domain"
)

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))
	claims := &jwt.RegisteredClaims{
		Subject:   user.Login,
		ExpiresAt: &jwt.NumericDate{Time: exp},
		ID:        fmt.Sprintf("%d", user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))
	claimsRefresh := &jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%d", user.ID),
		ExpiresAt: &jwt.NumericDate{Time: exp},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (int64, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return 0, fmt.Errorf("invalid Token")
	}
	jti, ok := claims["jti"].(string)
	if !ok {
		return 0, fmt.Errorf("jti claim is not a string")
	}
	id, err := strconv.ParseInt(jti, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse jti claim: %v", err)
	}
	return id, nil
}

func GetTokensByUser(user *domain.User, accessSecret string, refreshSecret string, accessExpiry int, refreshExpiry int) (AuthTokens, error) {
	accessToken, err := CreateAccessToken(user, accessSecret, accessExpiry)
	if err != nil {
		return AuthTokens{}, err
	}
	refreshToken, err := CreateRefreshToken(user, refreshSecret, refreshExpiry)
	if err != nil {
		return AuthTokens{}, err
	}
	authResp := AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return authResp, nil
}
