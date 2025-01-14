FROM node:20.18.1



RUN mkdir -p /opt/sowork
WORKDIR /opt/sowork
COPY . /opt/sowork
RUN npm install


WORKDIR /opt/sowork/src
EXPOSE 5000

CMD [ "node", "server.js" ]