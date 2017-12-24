from setuptools import setup, find_packages

if __name__ == "__main__":
    with open('README.md') as f:
        readme = f.read()

    setup(
        name='n0core',
        version='0.0.0',
        description='n0stack IaaS component',
        long_description=readme,
        url='https://github.com/n0stack/n0core',
        author='n0stack developer team',
        packages=find_packages()
    )
