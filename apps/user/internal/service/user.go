package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	v1 "vn.vato.zora.be.api/api/user/v1"
	"vn.vato.zora.be.api/apps/user/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServiceServer

	userUC  *biz.UserUseCase
	tokenUC *biz.TokenUseCase
	log     *log.Helper
}

func NewUserService(userUC *biz.UserUseCase, tokenUC *biz.TokenUseCase, logger log.Logger) *UserService {
	return &UserService{
		userUC:  userUC,
		tokenUC: tokenUC,
		log:     log.NewHelper(logger),
	}
}

func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	u, err := s.userUC.Register(ctx, req.Email, req.Password, req.Phone, req.Name)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterReply{UserId: u.ID.String(), Message: "registered successfully"}, nil
}

func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	userInfo, err := s.userUC.Login(ctx, req.Email, req.Password)
	if userInfo == nil || err != nil {
		return nil, err
	}

	access, refresh, expiresIn, err := s.tokenUC.GenerateTokenPair(ctx, userInfo.ID, userInfo.Role)
	if err != nil {
		return nil, err
	}

	return &v1.LoginReply{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         toProtoUser(userInfo),
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	access, expiresIn, err := s.tokenUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &v1.RefreshTokenReply{AccessToken: access, ExpiresIn: expiresIn}, nil
}

func (s *UserService) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutReply, error) {
	err := s.tokenUC.Logout(ctx, req.UserId, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return &v1.LogoutReply{Success: true}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.GetProfileReply, error) {
	u, err := s.userUC.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{User: toProtoUser(u)}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileReply, error) {
	u, err := s.userUC.UpdateProfile(ctx, req.UserId, req.Name, req.Phone, req.Avatar)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateProfileReply{User: toProtoUser(u)}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordReply, error) {
	if err := s.userUC.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}
	return &v1.ChangePasswordReply{Success: true, Message: "password changed"}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenReply, error) {
	userID, _, err := s.tokenUC.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &v1.ValidateTokenReply{Valid: false}, nil
	}
	return &v1.ValidateTokenReply{Valid: true, UserId: userID.String()}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *v1.GetUserByIdRequest) (*v1.GetUserByIdReply, error) {
	u, err := s.userUC.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserByIdReply{User: toProtoUser(u)}, nil
}

// ── helpers ─────────────────────────────────────────────

func toProtoUser(u *biz.UserInfo) *v1.User {
	return &v1.User{
		UserId:    u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
}
