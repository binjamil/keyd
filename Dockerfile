# Stage 1: Build source code into a statically linked binary
FROM golang:1.16 as build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o keyd

# Stage 2: Run the binary to start keyd server
FROM scratch
COPY --from=build /src/keyd .
EXPOSE 8080
CMD ["/keyd"]

# Copy .pem files and set env to enable TLS, for example:
# COPY --from=build /src/tls/*.pem .
# ENV TLS_CERTIFICATE="certificate.pem"
# ENV TLS_KEY="key.pem"