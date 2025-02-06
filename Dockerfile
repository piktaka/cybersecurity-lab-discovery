FROM node:20.18.1



RUN mkdir -p /opt/server
WORKDIR /opt/server
COPY . /opt/server

RUN npm install



EXPOSE 5000

CMD [ "node", "server.js" ]