FROM golang

WORKDIR /app
COPY main /app
COPY wait-for-it.sh /app

# CMD [ "./wait-for-it.sh", "rabbitmq:5672" ]
CMD [ "./main" ]
