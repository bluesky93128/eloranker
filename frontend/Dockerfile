FROM node:alpine

RUN apk add --no-cache git

WORKDIR /app

COPY package.json yarn.lock ./
# RUN yarn -s --no-cache --frozen-lockfile
RUN npm install -s --no-package-lock

COPY . ./

RUN yarn build

CMD rm -rf /dist/* && cp -r -f /app/dist/* /dist
