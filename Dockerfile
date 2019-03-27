FROM golang:1.12

COPY . .

RUN go build .

CMD ["isastudent"]