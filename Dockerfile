FROM golang:1.21.3-bookworm

# Setup mecab
RUN apt-get update -y \
	&& apt-get install -y mecab libmecab-dev mecab-ipadic-utf8

WORKDIR /livefetcher

COPY . .
RUN go mod download
RUN go build ./cmd/livefetcher

CMD ["make", "run"]