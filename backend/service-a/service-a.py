from flask import Flask, request
from confluent_kafka import Consumer, KafkaError, Producer
import threading
import os

app = Flask(__name__)

# Kafka configuration
KAFKA_BROKER = "kafka.event-queue.svc.cluster.local:9092"
SERVICE_NAME = "service-a"

consumer = Consumer({
    'bootstrap.servers': KAFKA_BROKER,
    'group.id': 'group-service-a',
    'auto.offset.reset': 'earliest'
})

producer = Producer({'bootstrap.servers': KAFKA_BROKER})


def consume_kafka_messages():
    consumer.subscribe([SERVICE_NAME])
    while True:
        msg = consumer.poll(1.0)
        if msg is None:
            continue
        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                continue
            else:
                print(msg.error())
                break
        else:
            correlation_id = msg.key().decode("utf-8")
            # 올바르게 문자열을 출력
            print(f"Received Kafka message: {msg.value().decode('utf-8')}, Correlation ID: {correlation_id}")


@app.route('/api/v1/service-a/test1', methods=['GET', 'POST'])
def http_check():
    correlation_id = request.headers.get('X-Correlation-ID', 'N/A')
    print(f"Received HTTP request with Correlation ID: {correlation_id}")
    return f"Service {SERVICE_NAME} HTTP Check! Correlation ID: {correlation_id}"


if __name__ == '__main__':
    threading.Thread(target=consume_kafka_messages).start()
    app.run(host='0.0.0.0', port=8081)
