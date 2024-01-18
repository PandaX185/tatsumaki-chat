FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o tatsumaki ./cmd/


# Expose the port your Go application will run on
EXPOSE 8080

# Command to run the application
CMD ["./tatsumaki"]
