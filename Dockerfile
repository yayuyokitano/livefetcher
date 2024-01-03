FROM golang:1.21.3-bookworm

# Setup mecab
RUN apt-get update -y \
	&& apt-get install -y libmecab-dev

ENV CGO_CFLAGS "-I/usr/include"
ENV CGO_LDFLAGS "-L/usr/lib/x86_64-linux-gnu -lmecab -lstdc++"

WORKDIR /livefetcher

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build

CMD ["make", "run"]