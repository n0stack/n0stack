import sys
sys.path.append('../../') # NoQA
from n0core.lib.messenger import Messenger
import pulsar


class MQHandler:
    handlers = dict()

    def __init__(self, service_url, topic, **kwargs):
        self.client = pulsar.Client(service_url)
        self.consumer = self.client.subscribe(topic, **kwargs)

    def handle(self, proto):
        def decorator(f):
            self._add_handler(proto, f)
            return f
        return decorator

    def _add_handler(self, proto, f):
        if proto in self.handlers:
            raise ValueError('{} handler already exists'.format(proto))
        self.handlers[proto] = f

    def listen(self):
        while True:
            inner_msg, messenger = Messenger.receive_message(self.consumer)
            protoname = inner_msg.DESCRIPTOR.name
            if protoname not in self.handlers:
                raise NotImplementedError('unhandled message: {}'.format(protoname))
            func = self.handlers[protoname]
            func(inner_msg, messenger)
