# Build stage
FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY . .
# Build the API application
RUN go build -o main main.go
# Must install cURL since its not included in alpine images
RUN apk add curl
# Download and install 'golang migrate' to run DB table scripts in Run stage
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

########################################################################################################################
# Run stage
FROM alpine:3.19
WORKDIR /app
# Copy the Executable that was created in the build stage after building the application
COPY --from=builder /app/main .
# Copy the 'golang migrate' tool we installed in the build stage
COPY --from=builder /app/migrate ./migrate
# Copy the API's env file
COPY app.env .
# Copy the migration scripts (table creations) so 'golang migrate' can run them
COPY db/migration ./migration
# Copy start.sh script we wrote. It will run 'golang migrate' and create all our tables in our database container (see start.sh for more)
COPY start.sh .
# Expose the port
EXPOSE 8080
# This is the argument being passed to the entry point command below
CMD ["/app/main"]
# Tell it to run start.sh script passing in "main.exe" as an argument
ENTRYPOINT ["/app/start.sh"]