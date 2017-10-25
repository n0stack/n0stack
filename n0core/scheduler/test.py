try:
    from n0core.lib import proto
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib import proto
try:
    from n0core.lib.n0mq import N0MQ
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ


client = N0MQ('pulsar://127.0.0.1:6650')
req = proto.CreateVMRequest(id='1', host='test')
producer = client.create_producer('persistent://main/sd/scheduler/handler')
producer.send(req)
