from ResorceCalculation import Execfunc


class Scheduler():
    """
    Scheduler main class
    """
    DB_USER = None
    DB_HOST = None
    DB_PORT = None
    DB_PASSWORD = None
    PROMETHEUS_HOST = None
    PROMETHEUS_PORT = None

    def __init__(self,
                 *,
                 db_user=None,  # type: Optional[str]
                 db_host=None,  # type: Optional[str]
                 db_port=None,  # type: Optional[str]
                 db_password=None,  # type: Optional[str]
                 prometheus_host=None,  # type: Optional[str]
                 prometheus_port=None,  # type: Optional[str]
                 ):
        self.DB_USER = db_user
        self.DB_HOST = db_host
        self.DB_PORT = db_port
        self.DB_PASSWORD = db_password
        self.PROMETHEUS_HOST = prometheus_host
        self.PROMETHEUS_PORT = prometheus_port

    def MQReceived(self, msg):
        if msg.host is None:
            self.scheduleJob()
        return self.send(proto, msg.host, msg.process)

    def ScheduleJob(self):
        re = Resorce()
        re.DBResorce(self.DB_USER, self.DB_PASSWORD, self.DB_HOST, self.DB_PASSWORD)
        re.PrometheusRresorce(self.PROMETHEUS_HOST, self.PROMETHEUS_PORT)
        return re.ResorceCalculation(Execfunc)

    def Send(proto, host, process):
        pass


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

    def DBResorce(user, password, host):
        pass

    def PrometheusRresorce():
        pass

    def ResorceCalculation(Execfunc):
        return Execfunc()
