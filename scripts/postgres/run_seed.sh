#!/bin/bash

echo "Creating Database tables"
psql -U $POSTGRES_USER -d $POSTGRES_DB -a -f /app/scripts/db/create_tables.sql

echo "Seeding tables with data"
psql -U $POSTGRES_USER -d $POSTGRES_DB -a -f /app/scripts/db/seed_data.sql