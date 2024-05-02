FROM eclipse-temurin:17-jdk-alpine as build
WORKDIR /staging/app

COPY mvnw .
COPY .mvn .mvn
COPY pom.xml .
COPY src src

RUN echo 'foo' > /tmp/foo.txt

RUN ./mvnw install -DskipTests
RUN mkdir -p target/dependency && (cd target/dependency; jar -xf ../*.jar)

FROM eclipse-temurin:17-jdk-alpine
VOLUME /tmp


ARG DEPENDENCY=/staging/app/target/dependency
COPY --from=build ${DEPENDENCY}/BOOT-INF/lib /app/lib
COPY --from=build ${DEPENDENCY}/META-INF /app/META-INF
COPY --from=build ${DEPENDENCY}/BOOT-INF/classes /app

WORKDIR /app

COPY ./Car_Inventory.csv .

ENTRYPOINT ["java","-cp",".:lib/*","com.ai.doomers.AidoomApplication"]
