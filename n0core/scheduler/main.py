import sys

from n0library.arguments.common import CommonArguments
from scheduler import Scheduler


sys.path.append('../../')


def main():
    argparser = CommonArguments(
        description="",
    ) # type: CommonArguments

    argparser.add_argument("--db-user",
                           type=str,
                           default=None,
                           dest='db_user',
                           help="")
    argparser.add_argument("--db-host",
                           type=str,
                           default=None,
                           dest='db_host',
                           help="")
    argparser.add_argument("--db-port",
                           type=int,
                           default=None,
                           dest='db_port',
                           help="")
    argparser.add_argument("--db-password",
                           type=str,
                           default=None,
                           dest='db_password',
                           help="")
    argparser.add_argument('--prometheus-host',
                           type=str,
                           default=None,
                           dest='prometheus_host',
                           help="")
    argparser.add_argument('--prometheus-port',
                           type=int,
                           default=None,
                           dest='prometheus_port',
                           help="")
    args = argparser.parse_args()


    scheduler = Scheduler(
        db_user=args.db_user,
        db_host=args.db_host,
        db_port=str(args.db_port),
        db_password=args.db_password,
        prometheus_host=args.prometheus_host,
        prometheus_port=str(args.prometheus_port),
    ) # type: Scheduler


if __name__ == '__main__':
    main()
