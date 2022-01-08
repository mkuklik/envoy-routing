FROM golang:alpine

WORKDIR /backend
COPY . /backend

RUN go build -o ./backend_main ./backend_main.go

RUN ls -l

CMD ["/backend/backend_main", "--logtostderr"]