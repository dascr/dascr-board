FROM node:14-alpine

ENV VITE_API_BASE=http://localhost:8000/
ENV VITE_API_URL=http://localhost:8000/api
ENV VITE_WS_URL=ws://localhost:8000/ws

WORKDIR /usr/src/app

COPY package*.json ./
COPY . .

RUN npm install
RUN npm run build

WORKDIR /usr/src/app/caddy

RUN wget "https://caddyserver.com/api/download?os=linux&arch=amd64" -O caddy
RUN chmod +x caddy

EXPOSE 5000

CMD [ "/usr/src/app/caddy/caddy", "run", "-config", "Caddyfile" ]