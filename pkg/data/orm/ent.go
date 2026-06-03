package orm

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"vn.vato.zora.be.api/pkg/logs"
)

type EntConfig struct {
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
	Dsn             string
	DbDriver        string
	Debugging       bool
}

func OpenConnection(ctx context.Context, logger log.Logger, conf EntConfig) dialect.Driver {
	driverType := conf.DbDriver
	l := logs.NewHelper(logger)
	if len(driverType) == 0 {
		driverType = defaultDriverType
	}

	var (
		db  *sql.DB
		err error
	)
	if conf.Debugging {
		cfg, err := pgx.ParseConfig(conf.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Tracer = tracer{l: l}

		connHandle := stdlib.RegisterConnConfig(cfg)
		db, err = sql.Open("pgx", connHandle)
		if err != nil {
			l.Fatal(ctx, err)
		}
	} else {
		db, err = sql.Open(driverType, conf.Dsn)
		if err != nil {
			l.Fatalf(ctx, "failed opening connection: %v", err)
		}
	}

	l.Debugf(ctx, "Connecting with driver: %s", driverType)
	l.Debugf(ctx, "DSN: %s", conf.Dsn)

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
		l.Errorf(ctx, "Database ping failed: %v, stats: %+v", err, stats)
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

	// subDir := filepath.Join("scripts", "db")
	// genScriptDir, err := utils.GetFirstExistingFile(subDir)
	// if err == nil {
	// 	err = executeSqlFile(c, drv, filepath.Join(genScriptDir, "immutable_unaccent.sql"))
	// 	if err != nil {
	// 		log.Warnf("failed executing immutable_unaccent.sql: %v", err)
	// 	}
	// }

	if conf.Debugging {
		sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
			// start := time.Now()
			// query := fmt.Sprint(i...)
			// elapsed := time.Since(start)
			// color := "\033[32m"
			// if elapsed > time.Second {
			// 	color = "\033[31m"
			// } else if elapsed > 200*time.Millisecond {
			// 	color = "\033[33m"
			// }
			// coloredTime := fmt.Sprintf("%s[GenTime: %s]%s | \033[34m[SQL]\033[0m %s", color, elapsed, "\033[0m", query)
			// l.WithContext(ctx).Debug(coloredTime)
			query := fmt.Sprint(i...)
			coloredTime := fmt.Sprintf("\033[34m[SQL]\033[0m %s", query)
			l.Debug(ctx, coloredTime)
		})
		return sqlDrv
	}

	return drv
}
