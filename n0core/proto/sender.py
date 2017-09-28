import all_pb2
import pulsar

client = pulsar.Client('pulsar://localhost:6650')
producer = client.create_producer(
    'persistent://sample/standalone/ns1/my-topic')

for i in range(10):
    producer.send('hello-pulsar-%d' % i)

client.close()
