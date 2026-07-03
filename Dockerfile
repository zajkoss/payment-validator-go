FROM golang:1.26-alpine as build
WORKDIR /app

COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o validator ./cmd/validator

FROM alpine AS runtime
WORKDIR /app
COPY --from=build /app/validator .
COPY testdata/ testdata/
ENTRYPOINT ["./validator"]