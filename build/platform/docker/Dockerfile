FROM golang:1.16-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /libgen-opds

FROM alpine:3.10
COPY --from=build /libgen-opds /libgen-opds
EXPOSE 5144
ENTRYPOINT ["/libgen-opds", "serve"]
