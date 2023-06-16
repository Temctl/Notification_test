FROM golang:1.16

# Set environment variables
ENV ENV_VARIABLE_NAME=value
ENV ANOTHER_VARIABLE_NAME=another_value

# Set the working directory
WORKDIR /app

# Copy the application source code
COPY . .

# Build the application
RUN go build -o myapp .

# Set the default command to run the application
CMD ["./myapp"]
