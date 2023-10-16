FROM golang:1.21

WORKDIR /usr/src/masonictemple4app


RUN mkdir -p /etc/env

# Note: Make sure you have the certs and creds folders in a local env folder
COPY . .

# It does not like that leading `.` so copy it
# manually, this works.
COPY env/.blog.env /etc/env/.blog.env

RUN mv env/* /etc/env/.

RUN go install

ENTRYPOINT ["masonictemple4.app", "runserver", "--port", "8080"]
