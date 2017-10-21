from typing import Any, Optional, Tuple, Dict # noqa

from base64 import b64encode, b64decode
from uuid import uuid4

import pulsar

from n0core.lib.proto import N0stackMessage


def generate_id():
    # type: () -> str
    return str(uuid4())


def parse_n0m(data):
    # type: (bytes) -> Any
    """Parse N0stackMessage and extract submessage

    Args:
        data: N0stackMessage to be processed

    Returns:
        Extracted submessage
    """
    msg = N0stackMessage()

    # TODO: Base64 solution should be replaced
    msg.ParseFromString(b64decode(data))

    msg_type = msg.WhichOneof('message')
    sub_msg = getattr(msg, msg_type)
    sub_msg_type = sub_msg.WhichOneof('message')

    return getattr(sub_msg, sub_msg_type)


def build_n0m(request_id, obj, type):
    # type: (str, Any, str) -> bytes
    """Construct N0stackMessage from submessage

    Args:
        request_id: Request ID for N0stackMessage
        obj: submessage payload
        type: type of N0stackMessage, 'Request' or 'Notification'
    """
    msg = N0stackMessage()
    msg.version = 1
    msg.request_id = request_id

    obj_type = obj.__class__.__name__
    getattr(getattr(msg, type), obj_type).MergeFrom(obj)

    # TODO: Base64 solution should be replaced
    return b64encode(msg.SerializeToString())


class N0MQProducer(pulsar.Producer):  # type: ignore
    def _build_msg(self, content, *args, **kwargs):  # type: ignore
        content = build_n0m(generate_id(), content, 'Request')
        return super()._build_msg(content, *args, **kwargs)


class N0MQConsumer(pulsar.Consumer):  # type: ignore
    def receive(self, *args, **kwargs):  # type: ignore
        msg = super().receive(*args, **kwargs)
        return parse_n0m(msg)


# override classes
pulsar.Producer = N0MQProducer
pulsar.Consumer = N0MQConsumer


class N0MQ(pulsar.Client):  # type: ignore
    """A class for handling N0stackMessage in the Pulsar MQ.

    N0stackMessage includes submessage Request or Notification,
    and N0MQ help us to extract submessage from N0stackMessage.

    Example:
        Initialization:
        >>> from n0core.lib.n0mq import N0MQ
        >>> from n0core.lib.proto import CreateVolumeRequest
        >>> client = N0MQ('pulsar://localhost:6650')

        Send new message with producer
        >>> producer = client.create_producer('persistent://sample/standalone/ns1/my-topic')
        >>> req = CreateVolumeRequest(id='1', host='test', size_mb=1)
        >>> producer.send(req)

        Set handler for receiving message with consumer
        >>> consumer = client.subscribe('persistent://sample/standalone/ns1/my-topic')
        >>> @consumer.on('CreateVolumeRequest')
        >>> def on_create_volume_request(message, auto_ack=False):
        >>>     print('CreateVolumeRequest')
        >>>     print(message.data)  # CreateVolumeRequest object
        >>>     consumer.ack(message)  # no need to ack if auto_ack is True (default)

        >>> consumer.listen()  # register all handler to message_listener of pulsar.Client.subscribe
    """

    handlers = dict()  # type: Dict

    def subscribe(self, topic, subscription_name=None):
        # type: (str, Optional[str]) -> N0MQHandler
        if subscription_name is None:
            subscription_name = generate_id()
        if (topic, subscription_name) in self.handlers:
            raise ValueError('subscriber for {}#{} already exists'.format(topic, subscription_name))
        handler = N0MQHandler(topic, subscription_name)
        self.handlers[(topic, subscription_name)] = handler
        return handler

    def do_subscribe(self, topic, subscription_name, *args, **kwargs):
        # type: (str, str, *str, **str) -> pulsar.Consumer
        return super().subscribe(topic, subscription_name, *args, **kwargs)

    def listen(self):
        # type: () -> None
        for ts in self.handlers:
            handler = self.handlers[ts]
            topic, subscription_name = ts
            consumer = self.do_subscribe(topic, subscription_name, message_listener=handler)
            handler.consumer = consumer
        while True:
            pass


class N0MQHandler(pulsar.Consumer):  # type: ignore
    handlers = dict()  # type: Dict

    def __init__(self, topic, subscription_name, *args, **kwargs):
        # type: (str, str, *str, **str) -> None
        self.topic = topic
        self.subscription_name = subscription_name

    def __call__(self, consumer, message):
        # type: (pulsar.Consumer, pulsar.Message) -> None
        message.data = parse_n0m(message.data())
        protoname = message.data.DESCRIPTOR.name
        if protoname not in self.handlers:
            raise NotImplementedError('unhandled message: {}'.format(protoname))
        func, auto_ack = self.handlers[protoname]
        func(message)
        if auto_ack:
            self.ack(message)

    def on(self, proto, auto_ack=True):
        # type: (str, bool) -> Any
        def wrapper(f):
            # type: (Any) -> Any
            self._add_handler(proto, f, auto_ack=auto_ack)
            return f
        return wrapper

    def _add_handler(self, proto, f, auto_ack):
        # type: (str, Any, bool) -> None
        if proto in self.handlers:
            raise ValueError('{} handler already exists on {}#{}'.format(
                    proto, self.topic, self.subscription_Name))
        self.handlers[proto] = f, auto_ack

    def ack(self, message):
        # type: (pulsar.Message) -> Any
        return self.consumer.acknowledge(message)
