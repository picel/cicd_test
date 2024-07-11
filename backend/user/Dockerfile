# Use an official Gradle image as the build stage
FROM gradle:latest AS build

WORKDIR /home/gradle/project

COPY --chown=gradle:gradle . .

RUN gradle build --no-daemon

# Use an official OpenJDK image as the base image for running the application
FROM openjdk:17-jdk-slim

WORKDIR /app

COPY --from=build /home/gradle/project/build/libs/*.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
