from n0library.logger import Logger
try:
    from n0core.lib.n0mq import N0MQ
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ


def send(url, req):
    producer = client.create_producer(url)
    producer.send(req)


client = N0MQ('pulsar://127.0.0.1:6650')
consumer = client.subscribe('persistent://main/sd/scheduler/handle')
logger = Logger(name='scheduler', stdout=True, level='info')
