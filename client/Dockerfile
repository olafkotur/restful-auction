FROM java:8-jdk-alpine

RUN mkdir /usr/app/
RUN mkdir /usr/app/client
RUN mkdir /usr/app/client/libs

COPY out/Client.class /usr/app/client

COPY libs/annotations-13.0.jar /usr/app/client/libs
COPY libs/kotlin-stdlib-1.3.50.jar /usr/app/client/libs
COPY libs/kotlin-stdlib-common-1.3.50.jar /usr/app/client/libs
COPY libs/okhttp-4.2.2.jar /usr/app/client/libs
COPY libs/okio-2.4.0.jar /usr/app/client/libs

WORKDIR /usr/app/client

CMD java -cp ./libs/annotations-13.0.jar:./libs/kotlin-stdlib-1.3.50.jar:./libs/kotlin-stdlib-common-1.3.50.jar:./libs/okhttp-4.2.2.jar:./libs/okio-2.4.0.jar:. Client