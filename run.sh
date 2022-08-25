#!/bin/zsh

# This is the bare minimum to run in development. For full list of flags,
# run ./vigilate -help
export $(grep -v '^#' .env | xargs -d '\n')

echo "> Building JS bundle file."
npm run build

echo
echo "> Building vigilate\n"
go build -o vigilate cmd/web/*.go && ./vigilate \
-db=$DEV_DB_NAME \
-domain=192.168.50.168 \
-dbuser=$DB_USER_ACC \
-dbpass=$DB_USER_PWD \
-pusherHost=$PUSHER_HOST \
-pusherKey=$PUSHER_KEY \
-pusherSecret=$PUSHER_SECRET \
-pusherApp=$PUSHER_APP \
-pusherPort=$PUSHER_PORT \
-pusherSecure=false

unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)