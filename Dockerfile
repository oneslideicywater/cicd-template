# Build the manager binary
FROM registry.geoway.com/golang/golang:1.18

WORKDIR /app

COPY . .
RUN chmod +x mc && /bin/cp -rf mc /bin/
# set up proxy
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
# clean dependency
RUN go mod tidy


# Build
# linux binary
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o cicd-template main.go
# windows binary
RUN CGO_ENABLED=0 GOOS=windows GO111MODULE=on GOARCH=amd64 go build -a -o cicd-template.exe main.go

RUN mc alias set minio/ http://172.16.66.37:9000 minio minio123 && mc cp cicd-template* minio/cicd-template/v1.0.0/
