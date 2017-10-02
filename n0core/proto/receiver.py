import n0stack_message_pb2 as message
import pulsar
import uuid
import time
import signal
import sys
import base64

client = pulsar.Client('pulsar://localhost:6650')
consumer = client.subscribe(
    'persistent://sample/standalone/ns1/my-topic',
    subscription_name='my-sub')


def signal_handler(signal, frame):
    client.close()
    sys.exit(0)


def parse_message(received):
    msg = message.N0stackMessage()
    msg.ParseFromString(base64.b64decode(received))

    msg_type = msg.WhichOneof('message')
    sub_msg = getattr(msg, msg_type)
    sub_msg_type = sub_msg.WhichOneof('message')

    return getattr(sub_msg, sub_msg_type)


def do_action(obj):
    print("type: %s" % obj.__class__.__name__)
    print(obj)

    # Do something by type
    if 'CreateVMRequest' == obj.__class__.__name__:
        pass
    elif 'UpdateVMRequest' == obj.__class__.__name__:
        pass
    elif 'DeleteVMRequest' == obj.__class__.__name__:
        pass
    elif 'UpdateVMPowerStateRequest' == obj.__class__.__name__:
        pass
    else:
        pass

signal.signal(signal.SIGINT, signal_handler)

while True:
    try:
        msg = consumer.receive(timeout_millis=10)
        # print("received")
        print("Received message: '%s'" % msg.data())
        # print(msg.data())
        do_action(parse_message(msg.data()))
        consumer.acknowledge(msg)
    except Exception as e:
        if e.args[0] != 'Pulsar error: TimeOut':
            print(e)

        pass

signal.pause()
client.close()
