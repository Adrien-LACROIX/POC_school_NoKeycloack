# Étape 1 : builder via Makefile 
FROM golang:1.23-alpine AS builder

WORKDIR /pocnokc

# Copie des sources et du Makefile
COPY . .

# Build via Makefile
RUN apk add --no-cache make && \
    make build-linux-amd64

# Étape 2 : image minimale
FROM alpine:latest

# Création d'un utilisateur non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copier le binaire depuis l'étape de build
COPY --from=builder /pocnokc/build/linux/amd64/POC_app-amd64 .

# Donner les droits à l'utilisateur
RUN chown appuser:appgroup ./POC_app-amd64

USER appuser

EXPOSE 8080

CMD ["./POC_app-amd64"]
