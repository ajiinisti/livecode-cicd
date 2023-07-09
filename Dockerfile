#STAGE 1
FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o livecode-cicd

#STAGE 2
from alpine
WORKDIR /app
COPY --from=build /src/livecode-cicd /app
ENTRYPOINT [ "/app/livecode-cicd" ]