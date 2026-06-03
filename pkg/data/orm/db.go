package orm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"entgo.io/ent/dialect"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/format"
	"vn.vato.zora.be.api/pkg/logs"
)

const (
	defaultSchemaName = "public"
	defaultDriverType = "pgx"
)

var mu sync.Mutex

type Connection struct {
	Driver dialect.Driver
	l      *logs.Helper
	name   string
}

func NewConnection(logger log.Logger, domainName, moduleName, dns, driver string, debugging bool,
) (*Connection, func(), error) {
	mu.Lock()
	defer mu.Unlock()
	var conn *Connection
	conn = new(Connection)
	conn.l = logs.NewHelper(logger)
	conn.name = fmt.Sprintf("%s-%s|%s", strings.ToLower(domainName), strings.ToLower(moduleName), format.TimeToString(time.Now(), constant.DEFAULT_DATE_FORMAT))

	// Open connection
	conn.openConnection(EntConfig{
		Dsn:       dns,
		DbDriver:  driver,
		Debugging: debugging,
	})

	return conn, conn.cleanup, nil
}

func (c *Connection) openConnection(conf EntConfig) {
	driverType := conf.DbDriver
	if len(driverType) == 0 {
		driverType = defaultDriverType
	}

	ctx := context.Background()
	var (
		db  *sql.DB
		err error
	)

	if conf.Debugging {
		cfg, err := pgx.ParseConfig(conf.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Tracer = tracer{l: c.l}

		connHandle := stdlib.RegisterConnConfig(cfg)
		db, err = sql.Open("pgx", connHandle)
		if err != nil {
			c.l.Fatal(ctx, err)
		}
	} else {
		db, err = sql.Open(driverType, conf.Dsn)
		if err != nil {
			c.l.Fatalf(ctx, "failed opening connection: %v", err)
		}
	}

	c.l.Debugf(ctx, "Connecting with driver: %s", driverType)
	c.l.Debugf(ctx, "DSN: %s", conf.Dsn)

	// TODO: change to dynamic driver type
	if driverType == defaultDriverType {
		driverType = dialect.Postgres
	}

	drv := entSql.OpenDB(driverType, db)
	if err != nil {
		panic(err)
	}

	err = drv.DB().Ping()
	if err != nil {
		stats := drv.DB().Stats()
		c.l.Errorf(ctx, "Database ping failed: %v, stats: %+v", err, stats)
		panic(err)
	}

	if conf.MaxOpenConn > 0 {
		drv.DB().SetMaxOpenConns(conf.MaxOpenConn)
	}
	if conf.MaxIdleConn > 0 {
		drv.DB().SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.ConnMaxLifetime > 0 {
		drv.DB().SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Minute)
	}

	if conf.Debugging {
		sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
			coloredTime := fmt.Sprint(i...)
			c.l.Debug(ctx, coloredTime)
		})
		c.Driver = sqlDrv
	} else {
		c.Driver = drv
	}
}

func (c *Connection) cleanup() {
	mu.Lock()
	defer mu.Unlock()
	c.l.Infof(context.Background(), "%s : closing the data resources", c.name)
	if err := c.Driver.Close(); err != nil {
		c.l.Error(context.Background(), err)
	}
	c.Driver = nil
}
