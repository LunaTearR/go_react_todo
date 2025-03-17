FROM node:20.13.1-alpine

WORKDIR /app

COPY package*.json /app/

RUN yarn

COPY . .

EXPOSE 3000

CMD [ "yarn", "dev" ]