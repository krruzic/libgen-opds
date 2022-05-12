FROM ubuntu
WORKDIR /app
RUN dpkg --add-architecture i386
RUN apt-get update
RUN apt-get install golang-go git libc6:i386 zlib1g:i386 -y
RUN git clone https://github.com/samsheff/Amazon-Kindle-Cross-Toolchain.git /opt
RUN apt-get install make -y
ENTRYPOINT ["/usr/bin/make"]
