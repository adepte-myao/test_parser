FROM golang:1.19-alpine AS build
WORKDIR /parser-service
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
WORKDIR /parser-service/cmd/app
RUN CGO_ENABLED=0 go build -o /bin/app

# Using alpine to have CLI, scratch dpesn't have bash
FROM alpine:latest 
RUN apk --update add ca-certificates
COPY --from=build /bin/app /bin/app
COPY --from=build /parser-service/config/config.yaml /config/config.yaml
ENV PATH=/bin

EXPOSE 9095
CMD [ "/bin/app" ]