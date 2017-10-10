FROM golang:1.9
RUN apt-get update && apt-get -y upgrade 

ADD ./ /go/src/github.com/kiwamunet/image-optim
RUN chmod +x /go/src/github.com/kiwamunet/image-optim/install.sh
RUN ["/bin/bash", "-c", "/go/src/github.com/kiwamunet/image-optim/install.sh"]

EXPOSE 8080
CMD ["image-optim"]