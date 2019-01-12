from setuptools import setup, find_packages

with open('README.md') as f:
    readme = f.read()

with open('LICENSE') as f:
    license = f.read()

with open('VERSION') as f:
    version = f.read()

setup(
    name='n0stack',
    version=version,
    # description='',
    long_description=readme,
    author='h-otter',
    author_email='h-otter@outlook.jp',
    install_requires=['protobuf', 'grpcio-tools', 'numpy'],
    url='https://github.com/n0stack/n0stack',
    license=license,
    packages=['n0test']+find_packages("n0proto.py"),
    package_dir={
        'n0test': 'build/n0test',
        'n0proto': 'n0proto.py',
    },
)
