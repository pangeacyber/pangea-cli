# Use the official Golang image as the base image
FROM golang:1.21@sha256:4746d26432a9117a5f58e95cb9f954ddf0de128e9d5816886514199316e4a2fb

RUN apt-get update && apt-get install -y sudo

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application files
COPY . .

# Download dependencies
RUN go mod download

# Run the install.sh script
RUN chmod +x ./install.sh

# Set the script as the entrypoint
ENTRYPOINT ["./install.sh"]
