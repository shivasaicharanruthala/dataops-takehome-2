# Use a minimal Ubuntu image for the second stage
FROM ubuntu:22.04

# Set the Current Working Directory inside the container
WORKDIR /dataops-takehome-2

# Copy the Pre-built binary file from the builder stage
COPY dataops-takehome .

# Command to run the executable
CMD ["./dataops-takehome"]
