package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 7 * 24 * time.Hour
	jwtSecret  = "your-secret-key"
)

type TokenRepo interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type TokenUseCase struct {
	repo TokenRepo
	log  *log.Helper
}

func NewTokenUseCase(repo TokenRepo, logger log.Logger) *TokenUseCase {
	return &TokenUseCase{repo: repo, log: log.NewHelper(logger)}
}

type RoleOA struct {
	OaId uuid.UUID `json:"oa_id"`
	Role []string  `json:"roles"`
}

type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	AccountType string    `json:"account_type"`
	RoleOAs     []RoleOA  `json:"role_oas"`
	jwt.RegisteredClaims
}

func (uc *TokenUseCase) GenerateTokenPair(ctx context.Context, userID uuid.UUID, roleOAs []RoleOA) (access, refresh string, expiresIn int64, err error) {
	claims := &Claims{
		UserID:      userID,
		AccountType: "",
		RoleOAs:     roleOAs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return
	}

	// refresh token
	refreshClaims := &Claims{
		UserID:  userID,
		RoleOAs: roleOAs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
		},
	}

	refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		return
	}

	// lưu refresh token vào Redis
	err = uc.repo.SaveRefreshToken(ctx, userID, refresh, refreshTTL)
	expiresIn = int64(accessTTL.Seconds())
	return
}

func (uc *TokenUseCase) ValidateToken(ctx context.Context, tokenStr string) (userID uuid.UUID, role []RoleOA, err error) {
	blacklisted, err := uc.repo.IsBlacklisted(ctx, tokenStr)
	if err != nil || blacklisted {
		return uuid.Nil, nil, status.Errorf(codes.Unauthenticated, "Token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, nil, status.Errorf(codes.Unauthenticated, "Token is invalid")
	}

	claims := token.Claims.(*Claims)
	return claims.UserID, claims.RoleOAs, nil
}

func (uc *TokenUseCase) RefreshToken(ctx context.Context, refreshToken string) (newAccess string, expiresIn int64, err error) {
	userID, err := uc.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", 0, status.Error(codes.Unauthenticated, "refresh token is invalid")
	}

	u := &Claims{UserID: userID}
	claims := &Claims{
		UserID: u.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
		},
	}

	newAccess, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	expiresIn = int64(accessTTL.Seconds())
	return
}

func (uc *TokenUseCase) Logout(ctx context.Context, userID, accessToken string) error {
	if err := uc.repo.BlacklistToken(ctx, accessToken, accessTTL); err != nil {
		return err
	}
	return uc.repo.DeleteRefreshToken(ctx, accessToken)
}
