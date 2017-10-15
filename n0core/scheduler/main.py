import sys
from urllib.parse import urlparse

from n0library.arguments.common import CommonArguments
from scheduler import Scheduler


sys.path.append('../../')


def main():
    argparser = CommonArguments(
        description="",
    )

    argparser.add_argument("--db-url",
                           type=str,
                           default=None,
                           dest="db_url",
                           help="")
    argparser.add_argument("--prometheus-url",
                           type=str,
                           default=None,
                           dest="prometheus_url",
                           help="")
    args = argparser.parse_args()

    db_url = urlparse(args.db_url)
    prometheus_url = urlparse(args.prometheus_url)

    scheduler = Scheduler(
        db_user=args.db_url.username,
        db_host=args.db_url.hostname,
        db_port=str(db_url.port),
        db_password=db_url.password,
        prometheus_host=prometheus_url.hostname,
        prometheus_port=str(prometheus_url.port),
    )


if __name__ == '__main__':
    main()
