env:
	docker-compose up

# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/esbuild:
	esbuild ./internal/templates/scripts/ --bundle --minify --outfile=./static/index.js --watch

live/templ:
	templ generate -lazy --watch --proxy="http://localhost:9000" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	@air \
	--build.cmd "GO_ENV=local go build -o bin/main cmd/api/*.go" --build.bin "bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# start live server and templ generation
live: 
	make -j3 live/server live/templ live/esbuild

compose-local:
	docker-compose -f docker-compose.local.yaml up

build:
	@GOOS=linux GOARCH=amd64 go build -o bin/main cmd/api/main.go

e2e:
	npm run test:ui

e2e:ci
	npm run test:ci

start:
	@GO_ENV=production ./bin/main

fetch-secrets:
	./deploy/scripts/fetch-secrets.sh ${ENV}

start-base: fetch-secrets
	ENV=${ENV} docker-compose -p base --env-file .env -f deploy/docker-compose.base.yaml up -d

stop-base:
	ENV=${ENV} docker-compose -p base --env-file .env -f deploy/docker-compose.base.yaml down

start-app: fetch-secrets
	ENV=${ENV} DOCKER_USERNAME=${DOCKER_USERNAME} DOCKER_TAG=${DOCKER_TAG} docker-compose -p app --env-file .env -f deploy/docker-compose.app.yaml up -d

stop-app: fetch-secrets
	ENV=${ENV} DOCKER_USERNAME=${DOCKER_USERNAME} DOCKER_TAG=${DOCKER_TAG} docker-compose -p app --env-file .env -f deploy/docker-compose.app.yaml down
