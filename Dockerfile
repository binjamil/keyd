# Stage 1: Build source code into a statically linked binary
FROM golang:1.16 as build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o keyd

# Stage 2: Run the binary to start keyd server
FROM scratch
COPY --from=build /src/keyd .
# ENV GRPC_ENABLED=true
# EXPOSE 50051
EXPOSE 8000
CMD ["/keyd"]