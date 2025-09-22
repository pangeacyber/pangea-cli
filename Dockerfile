# Use the official Golang image as the base image
FROM golang:1.25

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
