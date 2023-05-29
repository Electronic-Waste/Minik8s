FROM golang:latest

RUN apt update && apt install -y openssh-server sshpass openssh-client expect

COPY autossh.sh /go/src/

RUN /go/src/autossh.sh