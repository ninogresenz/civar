FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN env GOOS="${os}" GOARCH=amd64 CGO_ENABLED=0 go build -o civar --ldflags="-w -s"

FROM scratch
WORKDIR /app
COPY --from=builder /app/civar .
ENTRYPOINT ["./civar"]
CMD [""]