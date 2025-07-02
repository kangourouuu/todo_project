FROM golang:1.24.2-alpine AS BUILDER

WORKDIR /build

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main .

# FROM scratch 

# WORKDIR /app

# COPY --from=BUILDER /build /app/

EXPOSE 9003

CMD [ "/build/main" ]