package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vn.vato.zora.be.api/apps/user/internal/biz"
	"vn.vato.zora.be.api/apps/user/internal/data/ent"
	"vn.vato.zora.be.api/apps/user/internal/data/ent/accountlogin"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func (r UserRepo) FindAccountByEmail(ctx context.Context, email string) (*biz.Account, error) {
	result, err := r.data.Client.AccountLogin.
		Query().
		Where(accountlogin.Account(email)).
		WithUserAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find account: %v", err)
	}

	return toAccountBiz(result), nil
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r UserRepo) Create(ctx context.Context, u *biz.Account) error {
	current := time.Now()

	query := r.data.Client.AccountLogin.Create()
	query.SetAccount(u.Email).
		SetAccountType("USER").
		SetPasswordSalt(u.Salt).
		SetPasswordHash(u.Password).
		SetIsActive(true).
		SetCreatedAt(current).
		SetUpdatedAt(current).
		SetCreatedBy(u.CreatedBy).
		SetUpdatedBy(u.CreatedBy)

	login, err := query.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return status.Errorf(codes.AlreadyExists, "account already exists")
		}
		return status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	_, err = r.data.Client.UserAccount.Create().
		SetUserLoginID(login.ID).
		SetNillableFullName(&u.Name).
		SetNillablePhone(&u.Phone).
		SetCreatedAt(current).
		SetUpdatedAt(current).
		SetCreatedBy(u.CreatedBy).
		SetUpdatedBy(u.CreatedBy).
		Save(ctx)

	return nil
}

func (r UserRepo) ExistsAccountByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.data.Client.AccountLogin.
		Query().
		Where(accountlogin.Account(email)).
		Exist(ctx)
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to check email: %v", err)
	}
	return exists, nil
}

func (r UserRepo) FindByID(ctx context.Context, id string) (*biz.UserInfo, error) {
	// TODO implement me
	panic("implement me")
}

func (r UserRepo) Update(ctx context.Context, u *biz.UserInfo) error {
	// TODO implement me
	panic("implement me")
}

func (r UserRepo) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	// TODO implement me
	panic("implement me")
}

func (r UserRepo) AddAddress(ctx context.Context, a *biz.Address) error {
	// TODO implement me
	panic("implement me")
}

func (r UserRepo) ListAddresses(ctx context.Context, userID string) ([]*biz.Address, error) {
	// TODO implement me
	panic("implement me")
}

func (r UserRepo) DeleteAddress(ctx context.Context, userID, addressID string) error {
	// TODO implement me
	panic("implement me")
}

func toAccountBiz(l *ent.AccountLogin) *biz.Account {
	a := &biz.Account{
		ID:        l.ID,
		Email:     l.Account,
		Salt:      l.PasswordSalt,
		Password:  l.PasswordHash,
		CreatedAt: l.CreatedAt,
	}
	return a
}
