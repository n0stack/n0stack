import sys
sys.path.append('../../')
import pulsar

from n0core.lib.proto import UpdateVMPowerStateRequest
from n0core.lib.messenger import Messenger

def main():
    client = pulsar.Client('pulsar://localhost:6650')
    consumer = client.subscribe('persistent://sample/standalone/compute/114514',
                                subscription_name='compute')
    while True:

        inner_msg, messenger = Messenger.receive_message(consumer)
        print(inner_msg)
        print(messenger.request_id)
        
    client.close()


if __name__ == '__main__':
    main()
