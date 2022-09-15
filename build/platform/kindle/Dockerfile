FROM ubuntu:jammy
WORKDIR /app
RUN dpkg --add-architecture i386
RUN apt-get update
RUN apt-get install golang-go git libc6:i386 zlib1g:i386 -y
RUN git clone https://github.com/samsheff/Amazon-Kindle-Cross-Toolchain.git /opt

ENV CC=/opt/arm-kindle-linux-gnueabi/bin/arm-kindle-linux-gnueabi-gcc
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1
