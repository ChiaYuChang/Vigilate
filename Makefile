
include .env
DOMAIN := localhost

build-go:
	echo "> Building vigilate\n"
	go build -o vigilate cmd/web/*.go
	
build-webpack:
	echo "> Building JS bundle file."
	npm run build

build-all: build-webpack build-go

docker-new-psql-container:
	@docker volume create ${APP_NAME}-${APP_STATE}-postgres-volume
	@docker run --name ${APP_NAME}-${APP_STATE}-postgres \
	-p ${POSTGRES_PORT}:5432 \
	-v ${APP_NAME}-${APP_STATE}-postgres-volume:/var/lib/postgresql/data \
	-e POSTGRES_USER=${POSTGRES_USERNAME} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-e POSTGRES_DB=${POSTGRES_DB_NAME} \
	-d \
	${POSTGRES_IMAGE_TAG}

docker-create-db:
	@docker exec ${APP_NAME}-${APP_STATE}-postgres psql -U ${POSTGRES_USERNAME} -c 'CREATE DATABASE ${POSTGRES_DB_NAME};'

docker-flush-db:
	@docker container rm ${APP_NAME}-${APP_STATE}-postgres
	@docker volume rm ${APP_NAME}-${APP_STATE}-postgres-volume

docker-db-down:
	@docker stop ${APP_NAME}-${APP_STATE}-postgres

docker-db-up:
	@docker start ${APP_NAME}-${APP_STATE}-postgres

migrate-create:
	@echo "Name of .sql?: "; \
    read FILENAME; \
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq $${FILENAME} 

migrate-up:
	@docker run --rm -v ${MIGRATION_PATH}:/migrations --network host ${MIGRATION_IMAGE_TAG} -path /migrations  -database ${POSTGRESQL_URL} -verbose up 1

migrate-up-all:
	@docker run --rm -v ${MIGRATION_PATH}:/migrations --network host ${MIGRATION_IMAGE_TAG} -path /migrations  -database ${POSTGRESQL_URL} -verbose up

migrate-down:
	@docker run --rm -v ${MIGRATION_PATH}:/migrations --network host ${MIGRATION_IMAGE_TAG} -path /migrations  -database ${POSTGRESQL_URL} -verbose down 1

migrate-down-all:
	@docker run --rm -it -v ${MIGRATION_PATH}:/migrations --network host ${MIGRATION_IMAGE_TAG} -path /migrations  -database ${POSTGRESQL_URL} -verbose down

gen-ssl-certificate:
	@echo "Name of .crt and .key?: "; \
    read FILENAME; \
	openssl req -x509 -new -nodes -sha256 -utf8 \
	-days 3650 \
	-newkey rsa:2048 \
	-keyout ${KEY_PATH}/$${FILENAME}.key \
	-out ${KEY_PATH}/$${FILENAME}.crt \
	-config ssl.conf && \
	cat ${KEY_PATH}/$${FILENAME}.key ${KEY_PATH}/$${FILENAME}.crt > ${KEY_PATH}/$${FILENAME}.pem

gen-private-key:
	@openssl ecparam -genkey \
	-name secp384r1 \
	-out ${KEY_PATH}/${PRIVATE_KEY_NAME}

gen-public-key:
	@openssl req -new -x509 -sha256 \
	-key ${KEY_PATH}/${PRIVATE_KEY_NAME} \
	-out ${KEY_PATH}/${PUBLIC_KEY_NAME} \
	-config ssl.conf \
	-days 365

gen-ssl-key: gen-private-key gen-public-key 
	@cat ${KEY_PATH}/${PUBLIC_KEY_NAME} ${KEY_PATH}/${PRIVATE_KEY_NAME} > ${KEY_PATH}/${PEM_KEY_NAME} && \
	cp ${KEY_PATH}/${PUBLIC_KEY_NAME} ${CA_PATH}/

gen-pfx-file:
	openssl pkcs12 -export -in ${KEY_PATH}/${PUBLIC_KEY_NAME} -inkey ${KEY_PATH}/${PRIVATE_KEY_NAME} -out ${KEY_PATH}/${PFX_KEY_NAME}

gen-key: gen-private-key gen-public-key gen-pem-key
	@echo "Done"

run: build-all docker-db-up
	echo "> Start server"
	./vigilate \
		-db=${APP_NAME}-${APP_STATE} \
		-domain=${POSTGRES_HOST} \
		-dbuser=${POSTGRES_USERNAME} \
		-dbpass=${POSTGRES_PASSWORD} \
		-caDir=${CA_PATH} \
		-pusherHost=${PUSHER_HOST} \
		-pusherKey=${PUSHER_KEY} \
		-pusherSecret=${PUSHER_SECRET} \
		-pusherApp=${PUSHER_APP} \
		-pusherPort=${PUSHER_PORT} \
		-pusherSecure=false

clean:
	rm ./vigilate