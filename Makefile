
include .env
DOMAIN := 192.168.50.168

build-go:
	echo "> Building vigilate\n"
	go build -o vigilate cmd/web/*.go
	
build-webpack:
	echo "> Building JS bundle file."
	npm run build

build-all: build-webpack build-go

run: build-all
	echo "> Start server"
	./vigilate \
		-db=${DEV_DB_NAME} \
		-domain=${DOMAIN} \
		-dbuser=${DB_USER_ACC} \
		-dbpass=${DB_USER_PWD} \
		-pusherHost=${PUSHER_HOST} \
		-pusherKey=${PUSHER_KEY} \
		-pusherSecret=${PUSHER_SECRET} \
		-pusherApp=${PUSHER_APP} \
		-pusherPort=${PUSHER_PORT} \
		-pusherSecure=false

clean:
	rm ./vigilate