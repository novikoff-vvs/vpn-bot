FROM --platform=linux/amd64 golang:1.23-alpine3.21 as build

WORKDIR /app

COPY ./src/ ./
COPY ./src/go.mod ./
COPY ./configs/ ./configs/
COPY ./.env/ ./.env

RUN apk --no-cache --update add build-base

RUN go mod tidy

RUN go build  -o /app/dist/discord ./cmd/main.go

FROM --platform=linux/amd64 alpine

USER nobody

COPY --from=build --chown=nobody:nobody /app/dist /app
COPY  --from=build --chown=nobody:nobody /app/configs/ /app/configs/
COPY  --from=build --chown=nobody:nobody /app/.env /app/.env

WORKDIR /app

ENTRYPOINT ["/app/discord"]
