FROM gradle:7.5.1-jdk11 as builder
WORKDIR /ms-java-prometheus
COPY . /ms-java-prometheus/
RUN gradle assemble --no-daemon

FROM openjdk:11-slim
ENV APP_DIR /application
ENV APP_FILE server-uber-jar.jar

EXPOSE 8084
WORKDIR $APP_DIR
COPY --from=builder /ms-java-prometheus/build/libs/*-fat.jar $APP_DIR/$APP_FILE

ENTRYPOINT ["sh", "-c"]
CMD ["java -jar $APP_FILE"]