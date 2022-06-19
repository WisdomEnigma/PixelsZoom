FROM golang:1.19beta1-alpine3.15

ENV CGO_ENABLED=0

ENV PORT=3000

RUN mkdir /app

ADD . /app

WORKDIR /app

# set persist storage
ARG WORKDIR_SPACE=/app_data

RUN mkdir -p ${WORKDIR_SPACE}

ENV Mounted_Workspace=/app${WORKDIR_SPACE}/apps.txt

COPY go.mod go.sum ./

RUN go mod tidy

#RUN go mod vendor

RUN go mod download

# build app 
RUN go build -o pixelszoom 

# testing
RUN go test ./...


# publish app port
EXPOSE 3000

# peristance storage 
VOLUME [ ${WORKDIR_SPACE} ]

CMD ["/app/pixelszoom"]
