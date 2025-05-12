FROM golang:1.22-alpine

WORKDIR /app

# Копируем только go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходники
COPY . .

# Копируем фронтенд в нужное место
COPY cmd/gonews/webapp /app/webapp

# Сборка
RUN cd cmd/gonews && go build -o /app/gonewsaggregator .

CMD ["/app/gonewsaggregator"]
