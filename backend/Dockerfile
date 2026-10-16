FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o backend .

FROM alpine
COPY --from=build /app/backend /backend
CMD ["/backend"]
