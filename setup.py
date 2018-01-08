from setuptools import setup, find_packages

if __name__ == "__main__":
    with open('README.md') as f:
        readme = f.read()

    requires = [
        # common use
        "git+https://github.com/n0stack/n0library.git",

        # n0core.model.network, n0core.model.nic
        "netaddr",
        "pyroute2",
        "python-iptables",
        "git+https://github.com/larsks/python-netns.git",

        # n0core.target.vm
        # "libvirt-python"
    ]

    extras = {
        "test": [
            "flake8",
            "mypy",
            "lxml",
        ],
        "docs": [
            "Sphinx",
            "sphinx_rtd_theme",
            "recommonmark",
            "commonmark",
        ],
    }

    setup(
        name='n0core',
        version='0.0.4',
        description='n0stack IaaS component',
        long_description=readme,
        url='https://github.com/n0stack/n0core',
        author='n0stack developer team',
        packages=find_packages(),
        install_requires=list(filter(lambda p: "http" not in p, requires)),
        dependency_links=list(filter(lambda p: "http" in p, requires)),
        tests_requires=extras['test'],
        extras_require=extras,
    )
