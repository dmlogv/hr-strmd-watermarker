# Go builder
FROM golang:1.12.0 AS go_build

WORKDIR /app
COPY backend .

RUN go build -o watermarker .


# Vue builder
FROM node:8.15.1-alpine AS node_build

WORKDIR /app
COPY frontend .

RUN npm install \
    && npm run build


# Server
FROM scratch

WORKDIR /app
COPY --from=go_build   /app/watermarker .
COPY --from=node_build /app/frontend .

EXPOSE 3210:3210
CMD ["/app/watermarker -server 3210"]

