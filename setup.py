from setuptools import setup

with open('README.md') as f:
    readme = f.read()

with open('LICENSE') as f:
    license = f.read()

setup(
    name='n0stack',
    version='0.1.2',
    # description='',
    long_description=readme,
    author='h-otter',
    author_email='h-otter@outlook.jp',
    install_requires=['protobuf', 'grpcio-tools'],
    url='https://github.com/n0stack/n0stack',
    license=license,
    packages=['n0proto'],
)
