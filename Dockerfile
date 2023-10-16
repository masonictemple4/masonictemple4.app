FROM golang:1.21

WORKDIR /usr/src/masonictemple4app


RUN mkdir -p /etc/env
COPY . .

# It does not like that leading `.` so copy it
# manually, this works.
COPY env/.blog.env /etc/env/.blog.env

RUN mv env/* /etc/env/.

# RUN CGO_ENABLED=0 go build -o /go/bin/masonictemple4app main.go
RUN go install

# FROM scratch

# COPY --from=build /etc/env /etc/env
# # Just to be safe.
# COPY --from=build /etc/env/.blog.env /etc/env/.blog.env
# COPY --from=build /go/bin/masonictemple4app /bin/masonictemple4app
# 
# ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/env/creds/theswequarrycredentials.json


# ENTRYPOINT ["/bin/masonictemple4app", "runserver", "--port", "8080"]
ENTRYPOINT ["masonictemple4.app", "runserver", "--port", "8080"]
