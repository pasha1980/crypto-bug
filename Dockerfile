FROM golang:1.17 as cryptobug_builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN apt install -y tzdata
RUN CGO_ENABLED=0 go build -o main .
FROM alpine:3.15 as cryptobug_backend
COPY --from=cryptobug_builder /app/main /
COPY --from=cryptobug_builder /usr/share/zoneinfo /usr/share/zoneinfo
EXPOSE 80
CMD ["/main"]