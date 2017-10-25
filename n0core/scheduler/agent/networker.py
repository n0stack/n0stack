from initialize import consumer, logger, send  # NOQA
try:
    from n0core.lib import proto
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib import proto  # NOQA
