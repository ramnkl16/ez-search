
FROM golang:alpine AS builder
#RUN apk --no-cache add ca-certificates
# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Run test
#RUN go test ./...

# Build the application
RUN go build -o Product-int-logs .
#RUN apk add ca-certificates
# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to prerender folder
RUN cp /build/Product-int-logs .
#RUN cp /build/config.json .

############################
# STEP 2 build a small image
############################
# FROM chromedp/headless-shell:latest
# RUN apt install dumb-init
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /dist/Product-int-logs /
COPY ./config.json /config.json

# Command to run the executable
ENTRYPOINT ["/Product-int-logs", "-c", "config.json" -restapi]