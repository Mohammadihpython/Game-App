FROM ubuntu:latest
LABEL authors="hameddev"

ENTRYPOINT ["top", "-b"]