import sqlalchemy as sa
from urllib.parse import urlparse
from n0library.arguments.common import CommonArguments
from n0library.logger import Logger
try:
    from n0core.lib.n0mq import N0MQ
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

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


client = N0MQ(pulsar_url)
consumer = client.subscribe(scheduler_topic)
logger = Logger(name='scheduler', stdout=True, level='info')


def send(url, req):
    producer = client.create_producer(url)
    producer.send(req)


def CheckHost(msg):
    if msg.host is None:
        print("host")
    return send(msg)


class Scheduler():
    """
    Scheduler Method
    各protoタイプで実行するようなjobをメソッド化しておく
    """
    # postgresql+psycopg2://user:password@host:port/dbname[?key=value&key=value...]
    DB_URL = None
    PROMETHEUS_HOST = None
    PROMETHEUS_PORT = None

    def __init__(self,
                 *,
                 db_url=None,  # type: Optional[str]
                 prometheus_host=None,  # type: Optional[str]
                 prometheus_port=None,  # type: Optional[str]
                 ):
        self.DB_USER = db_url
        self.PROMETHEUS_HOST = prometheus_host
        self.PROMETHEUS_PORT = prometheus_port

    def ScheduleJob(self):
        re = Resorce()
        re.DBResorce(self.DB_URL)
        re.PrometheusRresorce(self.PROMETHEUS_HOST, self.PROMETHEUS_PORT)
        return re.ResorceCalculation(Execfunc)


scheduler = Scheduler(
    db_url=db_url,
    prometheus_host=prometheus_url.hostname,
    prometheus_port=str(prometheus_url.port),
)


class Resorce():
    """
    各リソースをclass変数でもっておき、メソッドが実行で取得し初期化する
    """
    DB_CPU = {}
    DB_MEMORY = {}
    DB_DISKFREE = {}
    Prometheus_CPU = {}
    Prometheus_MEMORY = {}
    Prometheus_DISKFREE = {}
    CREATE_VM_CPU = None
    CREATE_VM_MEMORY = None
    CREATE_VM_DISK = None

    def DBResorce(self, db_url):
        engine = sa.create_engine(db_url, echo=True)
        Session = sa.orm.sessionmaker(bind=engine)
        session = Session()
        hosts = session.query('HOST').all()
        for h in hosts:
            self.DB_CPU.append(h.cpu)
            self.DB_MEMORY.append(h.memory)
            self.DB_DISKFREE.append(h.diskfree)

    def PrometheusRresorce(self):
        return

    def ResorceCalculation(self, Execfunc):
        # Resorce 計算アルゴリズムは、後で自由に変えてくれという気持ち
        return Execfunc()
