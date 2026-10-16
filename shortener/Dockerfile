FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o shortener .

FROM alpine
COPY --from=build /app/shortener /shortener
CMD ["/shortener"]
