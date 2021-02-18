FROM golang:1.16-rc-buster
RUN mkdir /app
ADD . /app
WORKDIR /app
EXPOSE 8080
CMD ["go","run", "main.go"]