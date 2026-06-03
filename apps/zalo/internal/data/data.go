package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"vn.vato.zora.be.api/apps/zalo/internal/conf"
	"vn.vato.zora.be.api/apps/zalo/internal/data/ent"
	_ "vn.vato.zora.be.api/apps/zalo/internal/data/ent/runtime"
	"vn.vato.zora.be.api/apps/zalo/internal/data/ent_hook"
	"vn.vato.zora.be.api/apps/zalo/internal/data/ent_interceptor"
	"vn.vato.zora.be.api/pkg/data/orm"
	"vn.vato.zora.be.api/pkg/data/redis_db"
	"vn.vato.zora.be.api/pkg/logs"
)

var (
	domainName = "zora"
	moduleName = "zalo"
)

type Data struct {
	RedisClient *redis_db.RedisClient
	Conn        *orm.Connection
	Client      *ent.Client
	l           *logs.Helper
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := new(Data)
	data.l = logs.NewHelper(logger)
	var mode string
	if c.Database.Mode != nil {
		mode = strings.ToLower(*c.Database.Mode)
	}
	var schemaName string
	if c.Database.Schema != nil {
		schemaName = strings.ToLower(*c.Database.Schema)
	}

	dbConn, cleanupFunc, err := orm.NewConnection(
		logger,
		domainName,
		moduleName,
		c.Database.Source,
		c.Database.Driver,
		c.Database.Debugging,
	)
	if err != nil {
		panic(err)
	}

	data.Conn = dbConn
	data.newEnt()
	data.registerEntEvents()
	data.connectRedis(c.Redis)

	if c.Database.AutoMigrate {
		if err = data.AutoMigrate(context.Background(), mode, schemaName); err != nil {
			panic(err)
		}
	}

	return data, cleanupFunc, nil

}

func (d *Data) newEnt() {
	opts := []ent.Option{ent.Driver(d.Conn.Driver)}
	d.Client = ent.NewClient(opts...)
}

func (d *Data) connectRedis(c *conf.Data_Redis) {
	d.RedisClient = redis_db.NewRedisClient()
	config := map[string]any{
		"addr":     c.Address,
		"password": c.Password,
		"username": c.Username,
		"db":       c.Db,
	}
	err := d.RedisClient.Connect(config, func(err error) {
		if err != nil {
			d.l.Warnf(context.Background(), "Redis connect error: %v", err)
		}
	})
	if err != nil {
		panic(err)
	}
}

func (d *Data) registerEntEvents() {
	// Register ent hooks
	ent_hook.RegisterHooks(d.Client)
	// Register ent interceptors
	ent_interceptor.RegisterInterceptors(d.Client)
}

func (d *Data) WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, txErr := d.Client.Tx(ctx)
	if txErr != nil {
		return txErr
	}
	defer func() {
		if v := recover(); v != nil {
			if tx != nil {
				_ = tx.Rollback()
			}
			panic(v)
		}
	}()
	if fnErr := fn(tx); fnErr != nil {
		if rerr := tx.Rollback(); rerr != nil {
			fnErr = fmt.Errorf("%w: rolling back transaction: %v", fnErr, rerr)
		}
		return fnErr
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (d *Data) AutoMigrate(ctx context.Context, mode, schemaName string) error {
	opts := orm.GetMigrationOptions(mode, schemaName)
	d.l.Infof(ctx, "AutoMigrating ent_schema: %s", schemaName)
	err := d.Client.Schema.Create(ctx, opts...)
	if err != nil {
		d.l.Errorf(ctx, "AutoMigrate failed for ent_schema %s: %v", schemaName, err)
		return err
	}

	err = orm.ExecuteSqlScripts(ctx, d.Conn.Driver, schemaName, opts)
	if err != nil {
		d.l.Errorf(ctx, "Executing SQL scripts failed for ent_schema %s: %v", schemaName, err)
		return err
	}

	return nil
}
