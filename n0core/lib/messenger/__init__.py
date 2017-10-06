from base64 import b64encode, b64decode
from uuid import uuid4

import pulsar  # NOQA

from n0core.lib.proto import N0stackMessage

from typing import Any, Optional, Tuple  # NOQA


class Messenger(object):
    """
    Provides message queue sending function.

    Examples:
        Common initialization:
        >>> import pulsar
        >>> from n0core.lib.proto import CreateVMRequest
        >>> from n0core.lib.messenger import Messenger

        >>> client = pulsar.Client('pulsar://localhost:6650')

        >>> producer = client.create_producer(
        >>>     'persistent://sample/standalone/ns1/my-topic')

        Send new message:
        >>> req = CreateVMRequest(id="hogefuga")
        >>> Messenger.send_new_message(producer, req)

        Receive and sending new message having same request_id
        >>> consumer = client.subscribe(
        >>>     'persistent://sample/standalone/ns1/my-topic',
        >>>     subscription_name=str(uuid4()))

        >>> while True:
        >>>     inner_msg, messenger = Messenger.receive_message(consumer)
        >>>     # ... some action
        >>>     new_req = CreateVMRequest(id=str(uuid4()))
        >>>     messenger.send_message(producer, new_req)
    """

    def __init__(self, received_message=None):
        # type: (Optional[N0stackMessage]) -> None
        """
        Init messenger instance to create new messege.

        Args:
            received_message: Previous received messege having same request_id
        """
        if received_message:
            self.request_id = received_message.request_id
        else:
            self.request_id = Messenger.__generate_id()

    @classmethod
    def receive_message(cls, consumer):
        # type: (pulsar.Consumer) -> Tuple[Any, Messenger]
        """
        Receive message and generate Messenger instance.

        Args:
            consumer: consumer that receive message

        Returns:
            A Protobuf object defined as
            N0stackMessage.some_type.return_object_class
            and Messenger instance having request_id same to received
        """
        msg = consumer.receive()

        parsed_msg = cls.__parse_message(msg.data())

        instance = cls(received_message=parsed_msg)
        instance.parsed_msg = parsed_msg
        instance.inner_msg = cls.__get_inner_message(parsed_msg)

        consumer.acknowledge(msg)

        return instance.inner_msg, instance

    def receive_new_message(self, consumer):
        # type: (pulsar.Consumer) -> Any
        """
        Receive message and set self.parsed_msg, self.inner_msg.

        Args:
            consumer: consumer that receive message

        Returns:
            A Protobuf object defined as
            N0stackMessage.some_type.return_object_class
        """
        msg = consumer.receive()

        self.parsed_msg = Messenger.__parse_message(msg.data())
        self.inner_msg = Messenger.__get_inner_message(self.parsed_msg)
        self.request_id = self.parsed_msg.request_id

        consumer.acknowledge(msg)

        return self.inner_msg

    @classmethod
    def __parse_message(cls, received):
        # type: (str) -> N0stackMessage
        """
        Parse message received from message queue to Protobuf object.

        Args:
            received: a received message from pulsar

        Returns:
            Parsed message.
        """
        msg = N0stackMessage()

        # NOTE: Base64 decoding will be removed
        msg.ParseFromString(b64decode(received))

        return msg

    @classmethod
    def __get_inner_message(cls, msg):
        # type: (N0stackMessage) -> Any
        """
        Get inner message from N0stackMessage Protobuf object.

        Args:
            msg: A N0stackMessage Protobuf object.

        Returns:
            A Protobuf object defined as
            N0stackMessage.some_type.return_object_class
        """
        msg_type = msg.WhichOneof('message')
        sub_msg = getattr(msg, msg_type)
        sub_msg_type = sub_msg.WhichOneof('message')

        return getattr(sub_msg, sub_msg_type)

    @classmethod
    def send_new_message(cls, producer, obj, type='Request'):
        # type: (pulsar.Producer, Any, str, str) -> None
        """
        Send message to message queue with as a new request.

        Args:
            producer: a pulsar-client's Producer
            obj: Message payload.
              Must be defined as N0stackMessage.(type).(obj.__class__)
            type: Must be defined as N0stackMessage.(type)
        """
        request_id = cls.__generate_id()
        msg = cls.__construct_message(obj, request_id=request_id, type=type)
        producer.send(msg)

        return None

    def send_message(self, producer, obj, type='Request'):
        # type: (pulsar.Producer, Any, str, str) -> None
        """
        Send message to message queue with having request_id same to received.

        Args:
            producer: a pulsar-client's Producer
            obj: Message payload.
              Must be defined as N0stackMessage.(type).(obj.__class__)
            request_id: unique request id and identical to API request.
            type: Must be defined as N0stackMessage.(type)
        """
        msg = self.__class__.__construct_message(obj,
                                                 request_id=self.request_id,
                                                 type=type)
        producer.send(msg)

        return None

    @classmethod
    def __construct_message(cls, obj, request_id, type):
        # type: (Any, str, str) -> bytes
        """
        Serialize message from Protobuf object.

        Args:
            obj: Message payload.
              Must be defined as N0stackMessage.(type).(obj.__class__)
            request_id: unique request id and identical to API request.
            type: Must be defined as N0stackMessage.(type)

        Returns:
            A serialized message for Message Queue.
            (Currently, This encodes message with base64 due to restriction of
            pulsar-client, until supporting bytes receiver method)
        """
        msg = N0stackMessage()
        msg.version = 1
        msg.request_id = request_id

        obj_type = obj.__class__.__name__

        getattr(getattr(msg, type), obj_type).MergeFrom(obj)

        # NOTE: Base64 encoding will be removed
        return b64encode(msg.SerializeToString())

    @classmethod
    def __generate_id(cls):
        # type: () -> str
        return str(uuid4())
