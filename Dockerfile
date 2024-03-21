FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o user-service

FROM busybox:latest

WORKDIR /user-service

COPY --from=build /app/cmd/user-service .

COPY --from=build /app/cmd/.env .

EXPOSE 50001

CMD [ "./user-service" ]