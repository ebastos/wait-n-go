# Wait and Go

Based on other similar services, this can be used in conjuction with docker-compose to force a container to wait for certain services to be available before running it's main service.

## Example

docker-compose file:

```
version: "3.7"
services:
  hello:
    image: debian
    volumes:
      - ./wait-n-go:/bin/wait-n-go
    entrypoint: ["/bin/wait-n-go", "--services", "db:3306,rabbit:5672", "echo", "hello world"]
  db:
    image: mariadb:10.4.6
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: password
    ports:
        - 3306:3306
  rabbit:
    image: rabbitmq:3.7.16
    restart: always
    ports:
      - 5672:5672
      - 31501:5672
```

This will make sure the hello service will only start it's command after both MariaDB and RabbitMQ are online and listening.

## To build:

`go build .`

### Using Dagger

Make sure you have [Dagger](https://dagger.io/) installed.

Initialize your project:

```bash
dagger project init
dagger project update
```

Run the build:

```bash
dagger do build
```
