# Start from the latest golang base image
FROM golang:latest
# Set an image name
LABEL image="teltonika2tak-cot"

# Add Maintainer Info
LABEL maintainer="Rosen Velikov <rosen.st.velikov@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 7908

# Run the binary program produced by `go install`
CMD ["./main"]
