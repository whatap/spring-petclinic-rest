# Build Stage
FROM amazoncorretto:17-alpine AS build
WORKDIR /app

COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN chmod +x mvnw && ./mvnw dependency:resolve

COPY src ./src
RUN ./mvnw package

# Run Stage
FROM amazoncorretto:17-alpine
RUN apk update && \
  apk upgrade && \
  apk add --no-cache tar=1.35-r2 curl=8.12.1-r0  && \
  mkdir -p /data/whatap/bin /data/agent
COPY --from=build /app/target/*.jar /data/whatap/bin/app.jar
# 이미지 빌드하기 전 와탭 자바 에이전트를 ./whatap 경로에 다운로드 해둘 것
COPY whatap/whatap.agent-2.2.54.jar /data/agent
COPY entrypoint.sh /data/whatap/entrypoint.sh
ENTRYPOINT [ "/data/whatap/entrypoint.sh" ]
