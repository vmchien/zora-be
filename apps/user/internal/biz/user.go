package biz

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Account struct {
	ID        uuid.UUID
	Email     string
	Salt      string
	Password  string
	Name      string
	Phone     string
	Avatar    string
	Role      string
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	CreatedAt time.Time
}

type UserInfo struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Name      string
	Phone     string
	Avatar    string
	Role      []RoleOA
	CreatedAt time.Time
}

type Role struct {
	OaId uuid.UUID
	Role []string
}

type Address struct {
	ID        string
	UserID    string
	Label     string
	Street    string
	District  string
	City      string
	IsDefault bool
}

// Repository interface — data layer implement
type UserRepo interface {
	Create(ctx context.Context, u *Account) error
	ExistsAccountByEmail(ctx context.Context, email string) (bool, error)
	FindAccountByEmail(ctx context.Context, email string) (*Account, error)
	FindByID(ctx context.Context, id string) (*UserInfo, error)
	Update(ctx context.Context, u *UserInfo) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error

	AddAddress(ctx context.Context, a *Address) error
	ListAddresses(ctx context.Context, userID string) ([]*Address, error)
	DeleteAddress(ctx context.Context, userID, addressID string) error
}

type UserUseCase struct {
	userRepo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) Register(ctx context.Context, email, password, phone, name string) (*UserInfo, error) {
	existing, _ := uc.userRepo.ExistsAccountByEmail(ctx, email)
	if existing {
		return nil, status.Error(codes.InvalidArgument, "Email đã tồn tại")
	}

	salt := randomString(8)
	passwordHashed, err := hashPassword(password, salt)
	if err != nil {
		return nil, status.Error(codes.Internal, "Lỗi khi mã hóa mật khẩu")
	}

	u := &Account{
		Email:     email,
		Password:  passwordHashed,
		Salt:      salt,
		Phone:     phone,
		Name:      name,
		Role:      "CUSTOMER",
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, u); err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "Lỗi khi tạo tài khoản")
	}

	log.Infof("user registered: %s", u.ID)
	return &UserInfo{}, nil
}

func (uc *UserUseCase) ConfirmOTPRegister(ctx context.Context, email, password, phone, name string) (*UserInfo, error) {
	existing, _ := uc.userRepo.FindAccountByEmail(ctx, email)
	if existing != nil {
		return nil, status.Error(codes.InvalidArgument, "Email not found")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &Account{
		Email:     email,
		Password:  string(hashed),
		Phone:     phone,
		Name:      name,
		Role:      "CUSTOMER",
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	log.Infof("user registered: %s", u.ID)
	return &UserInfo{}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (*UserInfo, error) {
	account, err := uc.userRepo.FindAccountByEmail(ctx, email)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Tài khoản không hợp lệ!")
	}

	if err := comparePassword(password, account.Password, account.Salt); err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Sai mật khẩu")
	}

	return &UserInfo{}, nil
}

func (uc *UserUseCase) GetProfile(ctx context.Context, userID string) (*UserInfo, error) {
	return uc.userRepo.FindByID(ctx, userID)
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, userID, name, phone, avatar string) (*UserInfo, error) {
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "User not found")
	}

	if name != "" {
		u.Name = name
	}
	if phone != "" {
		u.Phone = phone
	}
	if avatar != "" {
		u.Avatar = avatar
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, userID, oldPass, newPass string) error {
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return status.Error(codes.FailedPrecondition, "User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPass)); err != nil {
		return status.Error(codes.FailedPrecondition, "Old password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdatePassword(ctx, userID, string(hashed))
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func hashPassword(password string, salt string) (string, error) {
	passwordSalt := salt + password
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwordSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func comparePassword(password string, passwordHash string, salt string) error {
	passwordSalt := salt + password
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordSalt))
}
