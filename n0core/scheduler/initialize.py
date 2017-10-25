import time
from _pulsar import ConsumerType
from urllib.parse import urlparse

from n0library.arguments.common import CommonArguments
from n0library.logger import Logger
try:
    from n0core.lib.n0mq import N0MQ
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

from scheduler import Scheduler
# from ResorceCalculation import Execfunc


parser = CommonArguments(
    description="",
)

parser.add_argument("--pulsar-url",
                    type=str,
                    default='pulsar://127.0.0.1:6650',
                    dest="pulsar_url",
                    help="Pulsar URL")
parser.add_argument("--scheduler-topic",
                    type=str,
                    default='persistent://main/sd/scheduler/handler',
                    dest="scheduler_topic",
                    help="Scheduler Topic")
parser.add_argument("--database-url",
                    type=str,
                    default=None,
                    dest="db_url",
                    help="Database URL")
parser.add_argument("--prometheus-url",
                    type=str,
                    default=None,
                    dest="prometheus_url",
                    help="Prometheus URL")
parser.add_argument("--compute-producer-url",
                    type=str,
                    default="persistent://sample/standalone/compute/",
                    dest="compute_producer",
                    help="Compute Producer URL")
parser.add_argument("--volumer-producer-url",
                    type=str,
                    default="persistent://sample/standalone/volumer/",
                    dest="volumer_producer",
                    help="Volumer Producer URL")
parser.add_argument("--porter-producer-url",
                    type=str,
                    default="persistent://sample/standalone/porter/",
                    dest="porter_producer",
                    help="Porter Producer URL")
parser.add_argument("--networker-producer-url",
                    type=str,
                    default="persistent://sample/standalone/networker/",
                    dest="networker_producer",
                    help="Networker Producer URL")
args = parser.parse_args()

pulsar_url = args.pulsar_url
scheduler_topic = args.scheduler_topic
db_url = urlparse(args.db_url)
prometheus_url = urlparse(args.prometheus_url)
compute_producer = args.compute_producer
volumer_producer = args.volumer_producer
porter_producer = args.porter_producer
networker_producer = args.networker_producer


class MQScheduler(N0MQ):

    def listen(self):
        # type: () -> None
        for ts in self.handlers:
            handler = self.handlers[ts]
            topic, subscription_name = ts
            consumer = self.do_subscribe(topic, subscription_name, message_listener=handler, consumer_type=ConsumerType.Exclusive)
            handler.consumer = consumer
        while True:
            time.sleep(100)


client = MQScheduler(pulsar_url)
consumer = client.subscribe(scheduler_topic)
logger = Logger(name='scheduler', stdout=True, level='info')

scheduler = Scheduler(
    db_url=db_url,
    prometheus_host=prometheus_url.hostname,
    prometheus_port=str(prometheus_url.port),
)


def send(url, req):
    producer = client.create_producer(url)
    producer.send(req)


def CheckHost(msg):
    if msg.host is None:
        print("host")
    return send(msg)
