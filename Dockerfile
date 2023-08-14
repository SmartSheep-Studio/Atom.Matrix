# Run image with this command
# docker run --rm --name matrix --net host -v $(pwd)/config.toml:/app/config.toml  -v $(pwd)/resources:/resources matrix

# Building Frontend
FROM node:18-alpine as matrix-web
WORKDIR /source
COPY . .
WORKDIR /source/packages/matrix-web
RUN rm -rf dist node_modules
RUN --mount=type=cache,target=/source/packages/matrix-web/node_modules,id=matrix_web_modules_cache,sharing=locked \
    --mount=type=cache,target=/root/.npm,id=matrix_web_node_cache \
    yarn install
RUN --mount=type=cache,target=/source/packages/matrix-web/node_modules,id=matrix_web_modules_cache,sharing=locked \
    yarn run build-only
RUN mv /source/packages/matrix-web/dist /dist

# Building Backend
FROM golang:alpine as matrix-server

WORKDIR /source
COPY . .
COPY --from=matrix-web /dist /source/packages/matrix-web/dist
RUN mkdir /dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /dist/server ./pkg/cmd/server/main.go

# Runtime
FROM golang:alpine

COPY --from=matrix-server /dist/server /app/server

EXPOSE 9443

CMD ["/app/server", "serve"]