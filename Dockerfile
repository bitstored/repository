FROM golang:alpine as source
WORKDIR /home/server
COPY . .
WORKDIR cmd/repository
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod vendor -o repository

FROM alpine as runner
LABEL name="bitstored/repository"
RUN apk --update add ca-certificates
COPY --from=source /home/server/cmd/repository/repository /bin/repository
EXPOSE 8080
ENTRYPOINT [ "repository" ]