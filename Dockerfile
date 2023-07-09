#STAGE 1
FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o simple-golang

#STAGE 2
from alpine
WORKDIR /app
COPY --from=build /src/simple-golang /app
ENTRYPOINT [ "/app/simple-golang" ]