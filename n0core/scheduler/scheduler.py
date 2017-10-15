import sqlalchemy as sa

from ResorceCalculation import Execfunc


class Scheduler():
    """
    Scheduler main class
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

    def MQReceived(self, msg):
        if msg.host is None:
            self.scheduleJob()
        return self.send(msg.host, msg.process)

    def ScheduleJob(self):
        re = Resorce()
        re.DBResorce(self.DB_URL)
        re.PrometheusRresorce(self.PROMETHEUS_HOST, self.PROMETHEUS_PORT)
        return re.ResorceCalculation(Execfunc)

    def Send(proto, host, process):
        return


class Resorce():
    """
    各リソースをclass変数でもっておき、メソッドが実行で取得し初期化する
    """
    DB_CPU = None
    DB_MEMORY = None
    DB_DISKFREE = None
    Prometheus_CPU = None
    Prometheus_MEMORY = None
    Prometheus_DISKFREE = None
    CREATE_VM_CPU = None
    CREATE_VM_MEMORY = None
    CREATE_VM_DISK = None

    def DBResorce(db_url):
        engine = sa.create_engine(db_url, echo=True)
        Session = sa.orm.sessionmaker(bind=engine)
        session = Session()
        hosts = session.query('HOST').all()
        return hosts

    def PrometheusRresorce():
        pass

    def ResorceCalculation(Execfunc):
        # Resorce 計算アルゴリズムは、後で自由に変えてくれという気持ち
        return Execfunc()
