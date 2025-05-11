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
	--build.cmd "go build -o tmp/bin/main cmd/api/*.go" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# start live server and templ generation
live: 
	make -j2 live/server live/templ