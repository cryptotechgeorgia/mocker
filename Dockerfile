FROM golang:alpine  as builder 

WORKDIR /app

COPY .  . 
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev
RUN go mod download

RUN go build -o mocker

FROM alpine:3.18.3 as production

RUN apk add --no-cache  tzdata


# Set the time zone
ENV TZ=Asia/Tbilisi
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ENV CGO_ENABLED=1


WORKDIR /app
COPY --from=builder /app/mocker .
ENV CGO_ENABLED=1

EXPOSE 8080

CMD ["./mocker"]
