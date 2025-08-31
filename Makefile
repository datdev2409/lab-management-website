env:
	docker-compose up

# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate -lazy --watch --proxy="http://localhost:8081" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	@air \
	--build.cmd "go build -o bin/main cmd/api/*.go" --build.bin "bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# start live server and templ generation
live: 
	make -j2 live/server live/templ

build:
	@GOOS=linux GOARCH=amd64 go build -o bin/main cmd/api/main.go

start:
	@GO_ENV=production ./bin/main

deploy:
	scp -i ~/Documents/AnhQuanLab/anhquanlab-mainserver.pem bin/main ec2-user@18.138.255.12:/home/ec2-user

copy-template:
	scp -i ~/Documents/AnhQuanLab/anhquanlab-mainserver.pem -r templates/ ec2-user@18.138.255.12:/home/ec2-user