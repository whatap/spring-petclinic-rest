build:
	./mvnw package
docker:
	docker build -t spring-petclinic-rest .
