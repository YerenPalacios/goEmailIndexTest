FROM golang:1.21.6
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
RUN apt update
RUN apt install -y graphviz

EXPOSE 8000

CMD ["/docker-gs-ping"]