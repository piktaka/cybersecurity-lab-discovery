FROM node:20.18.1



RUN mkdir -p /opt/sowork
WORKDIR /opt/sowork
COPY . /opt/sowork
RUN npm install


WORKDIR /opt/sowork/src
EXPOSE 5000
RUN apt update && apt install -y netcat-traditional
CMD [ "node", "server.js" ]