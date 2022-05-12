compile:
	CC=/opt/arm-kindle-linux-gnueabi/bin/arm-kindle-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build
