FROM node:15-alpine

ENV API_BASE=http://localhost:8000/
ENV API_URL=http://localhost:8000/api
ENV WS_URL=ws://localhost:8000/ws

WORKDIR /usr/src/app

COPY rollup.config.js ./
COPY package*.json ./
COPY . .

RUN npm install

RUN npm run build

EXPOSE 5000

CMD [ "npm", "run", "dev-expose" ]
