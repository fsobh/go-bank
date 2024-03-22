#!/bin/sh
# tell it to use bin/sh cuz alpine doesnt have bash

set -e # Telling it to exit immediately if a command returns a non 0 status

echo "run db migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up # Run the migrate up command

echo "start app"
exec "$@" # Execute the file path we passed in as an argument