# Set environment variables
$env:DB_HOST = "77.221.159.137"
$env:DB_PORT = "5432"
$env:DB_USER = "kratoff"
$env:DB_PASSWORD = "Oleg253535"
$env:DB_NAME = "encontro"
$env:DB_SSLMODE = "require"

# Load test data (assuming psql is in your PATH)
& psql "host=$env:DB_HOST port=$env:DB_PORT user=$env:DB_USER password=$env:DB_PASSWORD dbname=$env:DB_NAME sslmode=$env:DB_SSLMODE" -f "testdata.sql"

# Start the server
& go run "../../cmd/server/main.go"