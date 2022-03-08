FROM alpine:3.9

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY app .

EXPOSE 8080

CMD ["./app", "-C", "gs://${configBucket}/BitsyV1/config_v1.json"]
