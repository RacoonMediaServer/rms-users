FROM golang as builder
WORKDIR /src/service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-users -a -installsuffix cgo rms-users.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/rms-users .
COPY --from=builder /src/service/configs/rms-users.json /etc/rms/
CMD ["./rms-users"]