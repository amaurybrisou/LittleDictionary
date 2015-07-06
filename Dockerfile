FROM ubuntu:14.04
MAINTAINER Amaury Brisou

#ENTRYPOINT echo 'This is a GO Docker Image'

#install git
RUN apt-get update
RUN apt-get install -y curl git

# install go
RUN HOST_IP=$(netstat -nr | grep '^0\.0\.0\.0' | awk '{print $2}') && curl -s http://${HOST_IP}/repository/go1.4.2.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz
# install mongodb
RUN HOST_IP=$(netstat -nr | grep '^0\.0\.0\.0' | awk '{print $2}') && curl -s http://${HOST_IP}/repository/mongodb-linux-x86_64-ubuntu1404-3.0.4.tgz | tar -v -C /usr/local/ -xz

# setup go env
ENV PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:/usr/local/mongodb-linux-x86_64-ubuntu1404-3.0.4/bin
ENV GOPATH  /root/go
ENV GOROOT  /usr/local/go
ENV GOBIN   /usr/local/go/bin

RUN mkdir /root/db
RUN echo "mongod --dbpath /root/go/"

# ADD my sources to /root/go
ADD . /root/go/
# set working dir
WORKDIR /root/go
 	
# expose port 8080
EXPOSE 8080
EXPOSE 27017

#docker run -i -t -v `pwd`:/root/go -p 8080:8080 -p 27017:27017 -p 3000:3000 -p 9000:9000 amaurybrisou-go-aurita