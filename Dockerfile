# Go builder
FROM golang:1.12.0-alpine3.9 AS go_build


# Vue builder
FROM node:8.15.1-alpine AS node_build


WORKDIR /app
COPY frontend .
RUN npm install \
    && npm run build


# Server
FROM alpine:3.9

