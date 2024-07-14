# Database setup
cd internal/db/
sqlite3 home.db
.read home-init.sql
.exit

# Database dummy data setup
cd internal/db/
sqlite3 home.db
.read dummy-data.sql
.exit

# Database teardown
cd internal/db/
rm home.db