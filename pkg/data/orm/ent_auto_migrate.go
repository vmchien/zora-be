package orm

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"vn.vato.zora.be.api/pkg/utils"
)

type SchemaMigrator interface {
	SchemaCreate(ctx context.Context, schemaName string, opts ...schema.MigrateOption) error
}

type ModuleMigrator interface {
	SchemaCreate(ctx context.Context, domainName, moduleName, schemaName string, opts ...schema.MigrateOption) error
}

func AutoMigrateDB(client SchemaMigrator, conn dialect.Driver, mode string, customSchemaName string, sqlScriptFiles ...string) error {
	ctx := context.Background()

	// TODO: set false to drop index and column (MUST BE SET TO FALSE IN PRODUCTION)
	var migrationOptions []schema.MigrateOption
	if mode != "prod" {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(false),
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		}
	} else {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(false),
			schema.WithDropIndex(false),
			schema.WithDropColumn(false),
		}
	}

	if len(customSchemaName) > 0 {
		migrationOptions = append(migrationOptions, schema.WithSchemaName(customSchemaName))
	}

	// 🚀 Execute auto migration
	fmt.Println("🚀 Running auto migration...")
	// Execute any pre-migration SQL scripts if needed
	err := client.SchemaCreate(ctx, customSchemaName, migrationOptions...)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
		return err
	}
	// Execute additional SQL scripts if provided
	subDir := filepath.Join("internal", "data", "gen_sql")
	genScriptDir, err := utils.GetFirstExistingFile(subDir)
	if err != nil {
		log.Printf("⚠️ No SQL script files found in %s, skipping execution.", subDir)
		return nil
	}
	if len(sqlScriptFiles) > 0 {
		for _, fileName := range sqlScriptFiles {
			err = executeSqlFile(ctx, conn, filepath.Join(genScriptDir, fileName))
			if err != nil {
				log.Fatalf("❌ Migration failed: %v with error: %s", fileName, err.Error())
				return err
			}
		}
	}

	fmt.Println("✅ Migration completed successfully.")
	return nil
}

func AutoMigrateDBV2(client ModuleMigrator, conn dialect.Driver, mode, domainName, moduleName, schemaName string, sqlScriptFiles ...string) error {
	ctx := context.Background()

	// TODO: set false to drop index and column (MUST BE SET TO FALSE IN PRODUCTION)
	var migrationOptions []schema.MigrateOption
	if mode != "prod" {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(true),
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		}
	} else {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(true),
			schema.WithDropIndex(false),
			schema.WithDropColumn(false),
		}
	}

	if len(schemaName) > 0 {
		migrationOptions = append(migrationOptions, schema.WithSchemaName(schemaName))
	}

	// 🚀 Execute auto migration
	fmt.Println("🚀 Running auto migration...")
	// Execute any pre-migration SQL scripts if needed
	err := client.SchemaCreate(ctx, domainName, moduleName, schemaName, migrationOptions...)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
		return err
	}
	// Execute additional SQL scripts if provided
	subDir := filepath.Join("data", "gen_sql")
	genScriptDir, err := utils.GetFirstExistingFile(subDir)
	if err != nil {
		log.Printf("⚠️ No SQL script files found in %s, skipping execution.", subDir)
		return nil
	}
	if len(sqlScriptFiles) > 0 {
		for _, fileName := range sqlScriptFiles {
			err = executeSqlFile(ctx, conn, filepath.Join(genScriptDir, fileName))
			if err != nil {
				log.Fatalf("❌ Migration failed: %v with error: %s", fileName, err.Error())
				return err
			}
		}
	}

	fmt.Println("✅ Migration completed successfully.")
	return nil
}

func ExecuteSqlScripts(ctx context.Context, conn dialect.Driver, schemaName string, opts []schema.MigrateOption, sqlScriptFiles ...string) error {
	err := executeSystemFunctions(ctx, conn, schemaName)
	if err != nil {
		return err
	}

	// Execute additional SQL scripts if provided
	subDir := filepath.Join("internal", "data", "gen_sql")
	genScriptDir, err := utils.GetFirstExistingFile(subDir)
	if err != nil {
		fmt.Printf("⚠️ No SQL script files found in %s, skipping execution. \n", subDir)
		return nil
	}
	if len(sqlScriptFiles) > 0 {
		for _, fileName := range sqlScriptFiles {
			err = executeSqlFile(ctx, conn, filepath.Join(genScriptDir, fileName))
			if err != nil {
				log.Fatalf("❌ Migration failed: %v with error: %s", fileName, err.Error())
				return err
			}
		}
	}

	fmt.Println("✅ Migration completed successfully.")
	return nil
}

func GetMigrationOptions(mode, schemaName string) []schema.MigrateOption {
	var migrationOptions []schema.MigrateOption
	if mode != "prod" {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(false),
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		}
	} else {
		migrationOptions = []schema.MigrateOption{
			schema.WithGlobalUniqueID(false),
			schema.WithDropIndex(false),
			schema.WithDropColumn(false),
		}
	}

	if len(schemaName) > 0 {
		migrationOptions = append(migrationOptions, schema.WithSchemaName(schemaName))
	}
	return migrationOptions
}

// executeSqlFile reads and executes a SQL file using the given driver.
func executeSqlFile(ctx context.Context, drv dialect.Driver, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, skip execution
			return nil
		}
		return fmt.Errorf("failed to stat SQL file: %w", err)
	}
	if info.Size() == 0 {
		// Empty file, skip execution
		return nil
	}
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	statements := strings.Split(string(sqlBytes), ";")
	if len(statements) == 0 {
		return nil
	}

	tx, err := drv.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := tx.Exec(ctx, stmt, []any{}, nil); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute statement: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
