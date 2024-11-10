FROM golang:1.22.2 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM scratch AS final

EXPOSE 8080
COPY --from=build /app/main /app/main

ENTRYPOINT ["/app/main"]