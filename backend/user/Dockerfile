# 베이스 이미지로 OpenJDK 17 사용
FROM openjdk:17-jdk-slim AS build

# 작업 디렉토리 설정
WORKDIR /app

# gradle wrapper, build.gradle, settings.gradle 및 소스 코드 복사
COPY gradlew .
COPY gradle gradle
COPY build.gradle .
COPY settings.gradle .
COPY src src

# gradlew에 실행 권한 부여
RUN chmod +x gradlew

# Gradle 빌드 실행
RUN ./gradlew build --no-daemon

# 실행 단계 - JAR 파일만 포함하는 경량 이미지
FROM openjdk:17-jdk-slim

# 작업 디렉토리 설정
WORKDIR /app

# 빌드된 JAR 파일 복사
COPY --from=build /app/build/libs/*.jar app.jar

# 애플리케이션 실행
ENTRYPOINT ["java", "-jar", "app.jar"]
