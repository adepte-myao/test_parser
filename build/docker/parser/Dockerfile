FROM golang:1.19-alpine AS build
WORKDIR /parser-service
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
WORKDIR /parser-service/cmd/app
RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

# Using alpine to have CLI, scratch dpesn't have bash
FROM alpine:latest 
COPY --from=build /bin/app /bin/app
COPY --from=build /parser-service/config/config.yaml /config/config.yaml
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9095
CMD [ "/bin/app" ]