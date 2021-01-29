#!/bin/bash
set -euo pipefail

#psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
#	CREATE EXTENSION pg_stat_statements;
#EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "test_db" -c "CREATE EXTENSION pg_stat_statements;"
