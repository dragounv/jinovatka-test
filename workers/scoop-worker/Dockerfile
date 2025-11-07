FROM node:22.20 AS base

WORKDIR /app

COPY package.json ./

RUN npm install

RUN npx playwright install-deps chromium

COPY main.js ./

# Please mount apropriate config file in docker compose or on command line
#COPY config.json ./

ENTRYPOINT ["node", "main.js", "run"]