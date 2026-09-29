# syntax=docker/dockerfile:1
# Собирается из корня репозитория (нужен go.work):
#   docker build -f auth/Dockerfile -t auth .

FROM golang:1.25-alpine AS build
WORKDIR /src

# go.work требует, чтобы все модули из use(...) существовали, поэтому копируем все go.mod
COPY go.work go.work.sum ./
COPY auth/go.mod auth/go.sum ./auth/
COPY gateway/go.mod ./gateway/
COPY user/go.mod user/go.sum ./user/

COPY auth/ ./auth/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/auth ./auth/cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/auth /auth
EXPOSE 50051
ENTRYPOINT ["/auth"]
