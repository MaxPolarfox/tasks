FROM golang:1.13.1-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 3000

CMD ["./tasks_grpc"]