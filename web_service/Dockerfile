FROM node:6.2.0

EXPOSE 8004

COPY . /usr/src

WORKDIR /usr/src

RUN npm install

CMD node ./server/app.js
