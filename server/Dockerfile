FROM java:8-jdk-alpine

RUN mkdir /usr/app/
RUN mkdir /usr/app/server
RUN mkdir /usr/app/server/libs

COPY out/Server.class /usr/app/server

COPY libs/jedis-3.1.0.jar /usr/app/server/libs
COPY libs/json-20190722.jar /usr/app/server/libs

WORKDIR /usr/app/server
EXPOSE 8080

CMD java -cp ./libs/json-20190722.jar:./libs/jedis-3.1.0.jar:. Server