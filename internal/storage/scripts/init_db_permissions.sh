#!/bin/bash
# scripts/init_database.sh

# First check if .env exists
if [ ! -f ../../../.env ]; then
  echo "Error: .env file not found in parent directory"
  exit 1
fi

# Load ALL variables from .env file
# set -a allows the variables to be exported automatically
set -a
source ../../../.env
set +a

# Now we can create a temporary SQL file with our variables interpolated
cat >temp_setup.sql <<EOF
-- Create database if it doesn't exist
SELECT 'CREATE DATABASE ${DB_NAME}' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}')\gexec

-- Connect to the database
\c ${DB_NAME}

-- Create user if doesn't exist
DO
\$do\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '${DB_USER}') THEN
      CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
   END IF;
END
\$do\$;

-- Grant permissions
ALTER DATABASE ${DB_NAME} OWNER TO ${DB_USER};
ALTER SCHEMA public OWNER TO ${DB_USER};
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
GRANT ALL PRIVILEGES ON SCHEMA public TO ${DB_USER};
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${DB_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${DB_USER};
EOF

# Run the generated SQL file as postgres superuser
PGPASSWORD=${POSTGRES_PASSWORD} psql -U postgres -f temp_setup.sql

# Clean up the temporary file
rm temp_setup.sql

echo "Database setup complete! Database '${DB_NAME}' and user '${DB_USER}' have been created with appropriate permissions."
