FROM golang:1.16-alpine
## compile gcc statically
#ENV CGO_ENABLED=0
#ENV GOROOT=/usr/local/go
## this path will be mounted in deploy-service.yaml

#ENV PATH=$PATH:${GOROOT}/bin
#
## ATTENTION: you want to check, if the path to the project folder is the right one here
RUN mkdir /app
WORKDIR /app

COPY ./main /app/switcheroo

EXPOSE 9543

ENTRYPOINT ["./switcheroo"]
