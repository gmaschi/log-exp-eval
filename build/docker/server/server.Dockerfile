FROM golang:1.19.5 as builder

#Workdir
WORKDIR /go/src/github.com/gmaschi/log-exp-eval/

# Copying project to container filesystem
COPY . .

# Get packages required to run the application
RUN go mod tidy
RUN go mod download

# Compile the binary to execute the application
RUN go build -mod=mod -o main cmd/server/http/main.go

# Start the application
CMD ["./main"]
