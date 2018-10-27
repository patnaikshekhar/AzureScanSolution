FROM node:10.12.0-alpine
RUN npm install -g nodemon
ADD . /app
WORKDIR /app
RUN cd /app && npm install
EXPOSE 80
CMD ["nodemon", "index.js"]
