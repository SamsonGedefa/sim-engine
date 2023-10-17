FROM golang:alpine as dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go install github.com/githubnemo/CompileDaemon@latest

EXPOSE 8000
CMD CompileDaemon -polling -build="go build -o ./app" -command="./app"