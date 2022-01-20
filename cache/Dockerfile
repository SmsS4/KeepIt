FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY cache_runner ./cache_runner
COPY db ./
COPY ds ./
COPY utils ./
COPY *.yml ./
RUN go build -o /cache_runner