build_docker:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t libgen-opds:latest -f ./build/platform/docker/Dockerfile .

build_kual_extension:
	rm -rf ./build/dist/kindle
	mkdir -p ./build/dist/kindle
	cp -r ./build/platform/kindle/kual ./build/dist/kindle/libgen-opds
	docker build --platform=linux/amd64 -t libgen-opds-kindle-build -f build/platform/kindle/Dockerfile .
	docker run --rm -it -v `pwd`:/app libgen-opds-kindle-build go build -o ./build/dist/kindle/libgen-opds/bin/
	cd ./build/dist/kindle; zip -r "libgen-opds_`git describe --tags`.zip" -r ./libgen-opds; cd -
