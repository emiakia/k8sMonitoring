# Build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY k8s-mon-golang-app/go.mod k8s-mon-golang-app/go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the application source code
COPY k8s-mon-golang-app/*.go ./ 

# Build the application binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-monitoring-app

# Create the runtime container
FROM debian:bullseye-slim

# Install AWS CLI, necessary tools, and certificates
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    unzip \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip aws/ \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built application binary from the builder stage
COPY --from=builder /app/k8s-monitoring-app .

# Set the default KUBECONFIG environment variable
ENV KUBECONFIG=/root/.kube/config

# Command to run the application
CMD ["./k8s-monitoring-app"]

