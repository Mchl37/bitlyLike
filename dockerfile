FROM golang:1.21-bullseye
WORKDIR /app
COPY go.mod go.sum ./
COPY . /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
EXPOSE 1234
CMD ["/docker-gs-ping"]
