FROM golang:1.14.4-alpine3.12

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Daniil Gavrilovsky <d.gavrilovsky@corp.mail.ru>"

WORKDIR /tenjin

# Set depends
COPY go.mod go.sum ./
RUN go mod download
COPY . .

#create binary
RUN go build -o _build/userd ./userd/*.go
RUN chmod +x _build/userd

EXPOSE 8888

CMD ["_build/userd", "--config=_build/userd.yaml"]