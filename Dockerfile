FROM library/golang

# Godep for vendoring
RUN go get go.mongodb.org/mongo-driver
RUN go get github.com/beego/beego/v2@v2.0.0
RUN go get github.com/mxmCherry/translit
RUN go get golang.org/x/text/transform
RUN go get github.com/disintegration/imaging

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR $GOPATH/src/phonebook
RUN mkdir -p $APP_DIR

# Set the entrypoint
ENTRYPOINT (cd $APP_DIR && ./phonebook)
ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-d -w -s'

EXPOSE 8080
