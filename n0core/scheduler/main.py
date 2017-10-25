import pulsar  # NOQA
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
                    default='persistent://main/sd/scheduler/handle',
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
args = parser.parse_args()

pulsar_url = args.pulsar_url
scheduler_topic = args.scheduler_topic
db_url = urlparse(args.db_url)
prometheus_url = urlparse(args.prometheus_url)

scheduler = Scheduler(
    db_url=db_url,
    prometheus_host=prometheus_url.hostname,
    prometheus_port=str(prometheus_url.port),
)

client = N0MQ(pulsar_url)
consumer = client.subscribe(scheduler_topic)
logger = Logger(stdout=True, level='info')

from agent import compute  # NOQA
from agent import volumer, porter, networker  # NOQA


def send(url, req):
    producer = client.create_producer(url)
    producer.send(req)


if __name__ == '__main__':
    logger.info("listen start...")
    client.listen()
