# build frontend from source
FROM node:18-alpine AS build

ENV VITE_API_BASE=http://localhost:8000/
ENV VITE_API_URL=http://localhost:8000/api
ENV VITE_WS_URL=ws://localhost:8000/ws

WORKDIR /usr/src/app

COPY . .

RUN npm install; \
    npm run build


# create runtime image
FROM caddy:alpine

COPY ./caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /usr/src/app/build /usr/share/caddy

EXPOSE 8080
