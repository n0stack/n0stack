from typing import Any, Optional, Tuple # NoQA

from base64 import b64encode, b64decode
from uuid import uuid4

import pulsar

from n0core.lib.proto import N0stackMessage


def generate_id():
    return str(uuid4())

def parse_n0m(data):
    msg = N0stackMessage()
    
    # TODO: Base64 solution should be replaced
    msg.ParseFromString(b64decode(data))

    msg_type = msg.WhichOneof('message')
    sub_msg = getattr(msg, msg_type)
    sub_msg_type = sub_msg.WhichOneof('message')

    return getattr(sub_msg, sub_msg_type)

def build_n0m(request_id, obj, type):
    msg = N0stackMessage()
    msg.version = 1
    msg.request_id = request_id

    obj_type = obj.__class__.__name__
    getattr(getattr(msg, type), obj_type).MergeFrom(obj)

    # TODO: Base64 solution should be replaced
    return b64encode(msg.SerializeToString())


class N0MQProducer(pulsar.Producer):
    def _build_msg(self, content, *args, **kwargs):
        content = build_n0m(generate_id(), content, 'Request')
        return super()._build_msg(content, *args, **kwargs)


class N0MQConsumer(pulsar.Consumer):
    def receive(self, *args, **kwargs):
        msg = super().receive(*args, **kwargs)
        return parse_n0m(msg)


# override classes
pulsar.Producer = N0MQProducer
pulsar.Consumer = N0MQConsumer


class N0MQ(pulsar.Client):
    handlers = dict()

    def subscribe(self, topic, subscription_name=None):
        if subscription_name is None:
            subscription_name = generate_id()
        if (topic, subscription_name) in self.handlers:
            raise ValueError('subscriber for {}#{} already exists'.format(topic, subscription_name))
        handler = N0MQHandler(topic, subscription_name)
        self.handlers[(topic, subscription_name)] = handler
        return handler

    def do_subscribe(self, topic, subscription_name, *args, **kwargs):
        return super().subscribe(topic, subscription_name, *args, **kwargs)

    def listen(self):
        for ts in self.handlers:
            handler = self.handlers[ts]
            topic, subscription_name = ts
            consumer = self.do_subscribe(topic, subscription_name, message_listener=handler)
            handler.consumer = consumer
            while True:
                pass


class N0MQHandler(pulsar.Consumer):
    handlers = dict()

    def __init__(self, topic, subscription_name, *args, **kwargs):
        self.topic = topic
        self.subscription_name = subscription_name

    def __call__(self, consumer, message):
        message.data = parse_n0m(message.data())
        protoname = message.data.DESCRIPTOR.name
        if protoname not in self.handlers:
            raise NotImplementedError('unhandled message: {}'.format(protoname))
        func, auto_ack = self.handlers[protoname]
        func(message)
        if auto_ack:
            self.ack(message)

    def on(self, proto, auto_ack=True):
        def wrapper(f):
            self._add_handler(proto, f, auto_ack=auto_ack)
            return f
        return wrapper

    def _add_handler(self, proto, f, auto_ack):
        if proto in self.handlers:
            raise ValueError('{} handler already exists on {}#{}'.format(
                    proto, self.topic, self.subscription_Name))
        self.handlers[proto] = f, auto_ack

    def ack(self, message):
        return self.consumer.acknowledge(message)
