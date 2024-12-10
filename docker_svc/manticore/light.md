```shell
docker run -e EXTRA=1 --name manticore --rm -d manticoresearch/manticore && \
  until docker logs manticore 2>&1 | grep -q "accepting connections"; \
  do sleep 1; done && docker exec -it manticore mysql && docker stop manticore
```