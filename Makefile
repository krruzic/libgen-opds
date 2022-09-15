build_multiarch_docker:
	docker buildx build \
		--platform linux/amd64,linux/arm64,linux/arm/v7 \
		-t libgen-opds:latest \
		-f ./build/platform/docker/Dockerfile .

build_docker:
	docker build \
		-t libgen-opds:latest \
		-f ./build/platform/docker/Dockerfile .

build_multiarch:
	rm -rf ./build/dist/binary
	mkdir -p ./build/dist/binary
	go mod download
	env GOOS=darwin GOARCH=arm64 go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_darwin_arm64
	env GOOS=darwin GOARCH=amd64 go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_darwin_amd64
	env GOOS=linux GOARCH=amd64 go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_linux_amd64
	env GOOS=linux GOARCH=arm64 go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_linux_arm64
	env GOOS=linux GOARCH=arm go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_linux_arm
	env GOOS=windows GOARCH=amd64 go build -o ./build/dist/binary/libgen-opds_`git describe --tags`_windows_amd64.exe

build_kual_extension:
	rm -rf ./build/dist/kindle
	mkdir -p ./build/dist/kindle
	cp -r ./build/platform/kindle/kual ./build/dist/kindle/libgen-opds
	docker build --platform=linux/amd64 -t libgen-opds-kindle-build -f build/platform/kindle/Dockerfile .
	docker run --rm -it -v `pwd`:/app libgen-opds-kindle-build go build -o ./build/dist/kindle/libgen-opds/bin/
	cd ./build/dist/kindle; zip -r "libgen-opds_`git describe --tags`.zip" -r ./libgen-opds; cd -
