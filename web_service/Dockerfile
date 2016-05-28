FROM node:6.2.0

WORKDIR /usr/src

# copy package.json separately from code to optimize build cache for npm install
COPY package.json /usr/src
RUN npm install

COPY . /usr/src

EXPOSE 8080

CMD ["node", "./server/app.js"]

