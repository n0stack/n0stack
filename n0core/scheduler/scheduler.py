from typing import Any, Optional, Tuple, Dict  # NOQA
import sqlalchemy as sa

import Execfunc


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
        # type: () -> None
        self.DB_USER = db_url
        self.PROMETHEUS_HOST = prometheus_host
        self.PROMETHEUS_PORT = prometheus_port

    def ScheduleJob(self):
        # type: () -> str
        re = Resorce()
        re.DBResorce(self.DB_URL)
        re.PrometheusRresorce(self.PROMETHEUS_HOST, self.PROMETHEUS_PORT)
        return re.ResorceCalculation(Execfunc)


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

    def PrometheusResorce(self):
        return

    def ResorceCalculation(self, Execfunc):
        # Resorce 計算アルゴリズムは、後で自由に変えてくれという気持ち
        return Execfunc()
