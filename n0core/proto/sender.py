import n0stack_message_pb2 as message
import pulsar
import uuid
import time
import base64


def construct_message(obj, request_id='', type):
    msg = message.N0stackMessage()
    msg.version = 1
    msg.epoch_us = int(time.time() * 1e6)

    obj_type = obj.__class__.__name__

    getattr(getattr(msg, type), obj_type).MergeFrom(obj)

    return base64.b64encode(msg.SerializeToString())


def send_message(producer, obj, request_id='', type='Request'):
    msg = construct_message(obj, request_id=request_id, type=type)
    print(msg)
    producer.send(msg)

    return None

client = pulsar.Client('pulsar://localhost:6650')
producer = client.create_producer(
    'persistent://sample/standalone/ns1/my-topic')
