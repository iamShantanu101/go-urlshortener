## Running URL shortener with docker (single host)

### Prerequisites:
1. [Docker](https://www.docker.com/)

### Run the application with Docker:
1. Run `docker pull ishantanu16/gourlshortener:v1`
2. Run `docker container run -p 8080:8080 -d --name gourlshortener ishantanu16/gourlshortener:v1`
3. Run `docker ps | grep gourlshortener` to verify if the container is working fine

### Build the application locally and then run with Docker:
1. Run `docker build -t docker_id/app_name:[tag]` from the root directory to build the application locally
2. Run the application with `docker container run -p 8080:8080 -d --name gourlshortener docker_id/app_name:[tag]`
3. Run `docker ps | grep gourlshortener` to verify if the container is working fine
