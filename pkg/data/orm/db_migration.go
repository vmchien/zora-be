package orm

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
)

func (conn *Connection) SchemaCreate(ctx context.Context, schemaName string, opts ...schema.MigrateOption) error {
	if !strings.EqualFold(defaultSchemaName, schemaName) {
		err := conn.Driver.Exec(ctx, fmt.Sprintf("SET search_path TO %s", schemaName), []any{}, nil)
		if err != nil {
			conn.l.Errorf(ctx, "Failed to set search_path: %v", err)
			return err
		}
	}

	// Create ent_schema if it does not exist
	conn.l.Infof(ctx, "Creating ent_schema with options: %v", opts)

	// Migrate system function first if exists
	if err := executeSystemFunctions(ctx, conn.Driver, schemaName); err != nil {
		conn.l.Errorf(ctx, "Failed to execute system functions: %v", err)
		return err
	}

	err := migrateSchema(ctx, conn.Driver, opts...)
	if err != nil {
		conn.l.Errorf(ctx, "Failed to create ent_schema: %v", err)
		return err
	}
	conn.l.Infof(ctx, "Schema created successfully with options: %v", opts)
	return nil
}

func migrateSchema(ctx context.Context, drv dialect.Driver, opts ...schema.MigrateOption) error {
	m, err := schema.NewMigrate(drv, opts...)
	if err != nil {
		return err
	}
	return m.Create(ctx)
}

func migrateSchemaV2(ctx context.Context, drv dialect.Driver, opts ...schema.MigrateOption) error {

	m, err := schema.NewMigrate(drv, opts...)
	if err != nil {
		return err
	}
	return m.Create(ctx)
}

func executeSystemFunctions(ctx context.Context, drv dialect.Driver, schemaName string) error {
	scripts := []string{
		fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS unaccent"),
		fmt.Sprintf(`CREATE OR REPLACE FUNCTION immutable_unaccent(text)
							RETURNS text
							LANGUAGE sql IMMUTABLE PARALLEL SAFE AS
								  $func$
							SELECT %s.unaccent($1)
								  $func$;`, "public"),
	}

	for _, script := range scripts {
		if err := drv.Exec(ctx, script, []any{}, nil); err != nil {
			return fmt.Errorf("failed to execute system function script: %w", err)
		}
	}
	return nil

}

func (conn *Connection) ExecuteSqlScripts(ctx context.Context, schemaName string, opts ...schema.MigrateOption) error {
	if !strings.EqualFold(defaultSchemaName, schemaName) {
		err := conn.Driver.Exec(ctx, fmt.Sprintf("SET search_path TO %s", schemaName), []any{}, nil)
		if err != nil {
			conn.l.Errorf(ctx, "Failed to set search_path: %v", err)
			return err
		}
	}

	conn.l.Infof(ctx, "Creating ent_schema with options: %v", opts)

	// Migrate system function first if exists
	if err := executeSystemFunctions(ctx, conn.Driver, schemaName); err != nil {
		conn.l.Errorf(ctx, "Failed to execute system functions: %v", err)
		return err
	}

	err := migrateSchema(ctx, conn.Driver, opts...)
	if err != nil {
		conn.l.Errorf(ctx, "Failed to create ent_schema: %v", err)
		return err
	}
	conn.l.Infof(ctx, "Schema created successfully with options: %v", opts)
	return nil
}
