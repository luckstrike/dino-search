#!/bin/bash
# scripts/test_db_setup.sh

# Load environment variables
set -a
source ../../../../.env
set +a

# Function to check if the database exists
check_database() {
  echo "Checking if database '${DB_NAME}' exists..."
  if PGPASSWORD=${POSTGRES_PASSWORD} psql -U postgres -lqt | cut -d \| -f 1 | grep -qw ${DB_NAME}; then
    echo "✓ Database '${DB_NAME}' exists"
    return 0
  else
    echo "✗ Database '${DB_NAME}' does not exist"
    return 1
  fi
}

# Function to check if the user exists
check_user() {
  echo "Checking if user '${DB_USER}' exists..."
  if PGPASSWORD=${POSTGRES_PASSWORD} psql -U postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" | grep -q 1; then
    echo "✓ User '${DB_USER}' exists"
    return 0
  else
    echo "✗ User '${DB_USER}' does not exist"
    return 1
  fi
}

# Function to check if tables exist
check_tables() {
  echo "Checking if required tables exist..."
  local required_tables=("scraped_content" "content_headlines" "content_keywords")
  local all_exist=true

  for table in "${required_tables[@]}"; do
    if PGPASSWORD=${DB_PASSWORD} psql -U ${DB_USER} -d ${DB_NAME} -tAc "SELECT 1 FROM information_schema.tables WHERE table_name = '$table'" | grep -q 1; then
      echo "✓ Table '$table' exists"
    else
      echo "✗ Table '$table' does not exist"
      all_exist=false
    fi
  done

  $all_exist
  return $?
}

# Function to check if indexes exist
check_indexes() {
  echo "Checking if required indexes exist..."
  if PGPASSWORD=${DB_PASSWORD} psql -U ${DB_USER} -d ${DB_NAME} -tAc "SELECT 1 FROM pg_indexes WHERE indexname LIKE 'idx_scraped_content%'" | grep -q 1; then
    echo "✓ Required indexes exist"
    return 0
  else
    echo "✗ Required indexes are missing"
    return 1
  fi
}

# Main testing flow
echo "Starting database setup testing..."

# First, run the setup script
echo "Running database setup script..."
../init_db_permissions.sh

# Now run all our checks
failed=0

check_database || failed=1
check_user || failed=1
check_tables || failed=1
check_indexes || failed=1

# Try to connect and perform a test insert
echo "Testing data insertion..."
PGPASSWORD=${DB_PASSWORD} psql -U ${DB_USER} -d ${DB_NAME} <<EOF
INSERT INTO scraped_content (url, title, main_text) 
VALUES ('http://test.com', 'Test Title', 'Test Content') 
RETURNING id;
EOF

if [ $? -eq 0 ]; then
  echo "✓ Successfully inserted test data"
else
  echo "✗ Failed to insert test data"
  failed=1
fi

# Final result
if [ $failed -eq 0 ]; then
  echo -e "\n✅ All tests passed successfully!"
  exit 0
else
  echo -e "\n❌ Some tests failed. Please check the output above."
  exit 1
fi
