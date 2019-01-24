from setuptools import setup, find_packages


class Packages():
    def __init__(self):
        self._package = []
        self._package_dir = {}

    @property
    def package_dir(self):
        return self._package_dir

    @property
    def package(self):
        return self._package

    def add_package(self, package, directory):
        self._package_dir[package] = directory
        self._package.append(package)
        self._package.extend(self.__add_prefix(package, find_packages(directory)))

    @staticmethod
    def __add_prefix(prefix, l):
        return list(map(lambda x: prefix+"."+x, l))


if __name__ == "__main__":
    with open('README.md') as f:
        readme = f.read()

    with open('LICENSE') as f:
        license = f.read()

    with open('VERSION') as f:
        version = f.read()

    packages = Packages()
    packages.add_package('n0test', 'build/n0test')
    packages.add_package('n0proto', 'n0proto.py')

    print(packages.package)

    setup(
        name='n0stack',
        version="0.1."+version,
        description='A simple cloud provider using gRPC',
        long_description=readme,
        author='h-otter',
        author_email='h-otter@outlook.jp',
        install_requires=['protobuf', 'grpcio-tools', 'numpy'],
        url='https://github.com/n0stack/n0stack',
        license=license,
        packages=packages.package,
        package_dir=packages.package_dir,
    )
