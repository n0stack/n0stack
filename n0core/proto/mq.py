from base64 import b64encode, b64decode
import pulsar.Producer

import n0stack_message_pb2 as message

from typing import Any


def construct_message(obj, request_id, type):
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
    msg = message.N0stackMessage()
    msg.version = 1

    obj_type = obj.__class__.__name__

    getattr(getattr(msg, type), obj_type).MergeFrom(obj)

    # NOTE: Base64 encoding will be removed
    return b64encode(msg.SerializeToString())


def send_message(producer, obj, request_id='', type='Request'):
    # type: (pulsar.Producer, Any, str, str) -> None
    """
    Send message to message queue.

    Args:
        producer: a pulsar-client's Producer
        obj: Message payload.
          Must be defined as N0stackMessage.(type).(obj.__class__)
        request_id: unique request id and identical to API request.
        type: Must be defined as N0stackMessage.(type)
    """
    msg = construct_message(obj, request_id=request_id, type=type)
    print(msg)
    producer.send(msg)

    return None


def parse_message(received):
    # type: (str) -> message.N0stackMessage
    """
    Parse message received from message queue to Protobuf object.

    Args:
        received: a received message from pulsar

    Returns:
        Parsed message.
    """
    msg = message.N0stackMessage()
    # NOTE: Base64 decoding will be removed
    msg.ParseFromString(b64decode(received))

    return msg


def get_inner_message(msg):
    # type: (message.N0stackMessage) -> Any
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


# Common
# import pulsar
# client = pulsar.Client('pulsar://localhost:6650')

# Producer
# producer = client.create_producer(
#     'persistent://sample/standalone/ns1/my-topic')
# req = message.CreateVMRequest(id="hogevm", host="node01", ...)
# sender.send_message(producer, obj=req)


# Subscriber
# consumer = client.subscribe(
#     'persistent://sample/standalone/ns1/my-topic',
#     subscription_name='my-sub')
# try:
#     msg = consumer.receive(timeout_millis=10)
#     parsed_msg = parse_message(msg.data()))
#     # ... do some action corresponds with the type of message
#     consumer.acknowledge(msg)
# except Exception as e:
#     if e.args[0] != 'Pulsar error: TimeOut':
#         print(e)
#     pass
