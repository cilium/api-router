FROM ubuntu
WORKDIR /
ADD . /
ENTRYPOINT ["/api-router"]
