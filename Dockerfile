# Start from a Go image
FROM golang:1.24.3-alpine

# Setup Work Directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o golangcompoundtag .

# Expose the port your app will run on
EXPOSE ${PORT}

ENV DB_USER=sa
ENV DB_PASS=p$th@2024**
ENV DB_SERVER=PSTH-SRRYAPP04
ENV DB_NAME=DB_COMPOUND
ENV PORT=3420


# Run the application
CMD ["./golangcompoundtag"]
