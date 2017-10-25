from initialize import client, logger  # NOQA
from agent import compute, volumer, porter, networker  # NOQA


if __name__ == '__main__':
    logger.info("listen start...")
    client.listen()
