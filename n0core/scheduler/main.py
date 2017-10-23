from urllib.parse import urlparse

from n0library.arguments.common import CommonArguments
from scheduler import Scheduler


def main():
    parser = CommonArguments(
        description="",
    )

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

    db_url = urlparse(args.db_url)
    prometheus_url = urlparse(args.prometheus_url)

    scheduler = Scheduler(
        db_user=db_url.username,
        db_host=db_url.hostname,
        db_port=str(db_url.port),
        db_password=db_url.password,
        prometheus_host=prometheus_url.hostname,
        prometheus_port=str(prometheus_url.port),
    )


if __name__ == '__main__':
    main()
