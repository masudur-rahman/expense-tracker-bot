FROM ghcr.io/masudur-rahman/golang:1.20

WORKDIR /expense-tracker

COPY . .
#RUN go mod tidy && go mod vendor
RUN go build -o expense-tracker

#USER nobody:nobody
USER 65535:65535

EXPOSE 8080

ENTRYPOINT ["./expense-tracker"]
CMD ["serve"]
