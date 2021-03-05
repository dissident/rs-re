FROM rs-re-base AS build
WORKDIR /src
COPY . .
# RUN go build -o /out/example .

RUN go run *.go

# CMD ["/out/example"]
