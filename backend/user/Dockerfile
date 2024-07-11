# Use an official Gradle image to build the application
FROM gradle:7.4.2-jdk17 AS build

# Set the working directory in the container
WORKDIR /home/gradle/project

# Copy the project files to the container
COPY --chown=gradle:gradle . .

# Build the project
RUN gradle build --no-daemon

# Use an official OpenJDK image as the base image for running the application
FROM openjdk:17-jdk-slim

# Set the working directory in the container
WORKDIR /app

# Copy the JAR file from the build stage
COPY --from=build /home/gradle/project/build/libs/*.jar app.jar

# Expose the port the application runs on
EXPOSE 8080

# Define the command to run the application
ENTRYPOINT ["java", "-jar", "app.jar"]
