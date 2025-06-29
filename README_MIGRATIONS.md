# Database Migrations

This project includes a robust database migration system that helps manage database schema changes over time.

## Overview

The migration system provides:
- **Automatic migrations**: Runs on application startup
- **Manual migration control**: CLI tool for advanced operations
- **Version tracking**: Tracks applied migrations in the database
- **Rollback support**: Can rollback the last applied migration
- **Migration templates**: Automatic file generation with proper structure

## Migration File Format

Migration files follow the naming convention: `YYYYMMDDHHMMSS_migration_name.sql`

Each migration file contains both UP and DOWN sections:

```sql
-- +migrate Up
-- Your forward migration SQL here
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +migrate Down
-- Your rollback migration SQL here
DROP TABLE IF EXISTS example;
```

## Usage

### Automatic Migrations

Migrations run automatically when the application starts. The system:
1. Creates a `schema_migrations` table to track applied migrations
2. Scans the `migrations/` directory for migration files
3. Applies any pending migrations in order
4. Logs the progress

### Manual Migration Commands

Use the migration CLI tool for manual control:

```bash
# Run all pending migrations
make migrate-up

# Rollback the last applied migration
make migrate-down

# Check migration status
make migrate-status

# Create a new migration file
make migrate-create NAME=add_user_table
```

### Direct CLI Usage

You can also use the migration binary directly:

```bash
# Build the migration tool
make build-migrate

# Run migrations
./migrate -action=up

# Rollback last migration
./migrate -action=down

# Check status
./migrate -action=status

# Create new migration
./migrate -action=create -name=your_migration_name
```

## Migration Directory Structure

```
migrations/
├── 20241216120000_initial_schema.sql
├── 20241216120001_add_user_table.sql
└── 20241216120002_add_indexes.sql
```

## Best Practices

### 1. **Always provide rollback scripts**
Every migration should have a proper DOWN section that can cleanly reverse the changes.

### 2. **Test migrations thoroughly**
- Test both UP and DOWN migrations
- Test on a copy of production data
- Verify data integrity after migration

### 3. **Keep migrations atomic**
Each migration should focus on a single logical change:
- ✅ Good: `add_user_email_column.sql`
- ❌ Bad: `add_multiple_tables_and_indexes.sql`

### 4. **Use descriptive names**
Migration names should clearly describe what they do:
- ✅ Good: `add_user_authentication_table.sql`
- ❌ Bad: `update_schema.sql`

### 5. **Handle data migrations carefully**
When migrating data:
- Use transactions
- Consider performance impact
- Have a rollback plan for data changes

### 6. **Use IF EXISTS/IF NOT EXISTS**
Make migrations idempotent when possible:
```sql
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS users;
```

## Configuration

The migration system uses the same database configuration as the main application:
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

## Migration Table

The system creates a `schema_migrations` table to track applied migrations:

```sql
CREATE TABLE schema_migrations (
    version BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Troubleshooting

### Migration Failed
If a migration fails:
1. Check the error message in the logs
2. Fix the SQL in the migration file
3. The failed migration won't be marked as applied
4. Run migrations again after fixing

### Need to Skip a Migration
If you need to manually mark a migration as applied:
```sql
INSERT INTO schema_migrations (version, name) 
VALUES (20241216120001, 'migration_name');
```

### Rollback Multiple Migrations
Currently, rollback only supports rolling back one migration at a time. To rollback multiple migrations, run `make migrate-down` multiple times.

## Docker Integration

The migration system works seamlessly with Docker:
- Migrations run automatically when the app container starts
- The `migrations/` directory should be mounted or included in the Docker image
- Database connection uses the same environment variables as configured in `docker-compose.yml`

## Example Workflow

1. **Create a new migration**:
   ```bash
   make migrate-create NAME=add_user_profiles
   ```

2. **Edit the generated file** (`migrations/XXXXXX_add_user_profiles.sql`):
   ```sql
   -- +migrate Up
   CREATE TABLE user_profiles (
       id SERIAL PRIMARY KEY,
       user_id INTEGER REFERENCES users(id),
       bio TEXT,
       avatar_url VARCHAR(500),
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );

   -- +migrate Down
   DROP TABLE IF EXISTS user_profiles;
   ```

3. **Test the migration**:
   ```bash
   make migrate-status  # Check current status
   make migrate-up      # Apply the migration
   make migrate-down    # Test rollback
   make migrate-up      # Apply again
   ```

4. **Deploy**: The migration will run automatically when the application starts in production.

## Advanced Usage

### Custom Migration Path
You can customize the migrations directory by modifying the `migration.NewMigrator()` call in the datastore package.

### Multiple Environments
Use different migration directories or databases for different environments by setting appropriate environment variables.

### Seed Data
For seed data, create migrations that use `INSERT ... ON CONFLICT DO NOTHING` to ensure idempotency. 