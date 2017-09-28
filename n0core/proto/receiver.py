import all_pb2
import pulsar
import signal
import sys


client = pulsar.Client('pulsar://localhost:6650')
consumer = client.subscribe(
    'persistent://sample/standalone/ns1/my-topic',
    subscription_name='my-sub')


def signal_handler(signal, frame):
    client.close()
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

while True:
    try:
      msg = consumer.receive(timeout_millis=10)
      print("Received message: '%s'" % msg.data())
      consumer.acknowledge(msg)
    except Exception as e:
      pass


signal.pause()
client.close()
