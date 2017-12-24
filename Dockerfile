FROM python:3.6

COPY . /srv
WORKDIR /srv
RUN apt-get update \
 && apt-get install -y iptables \
 && pip install --no-cache-dir \
        pulsar-client \
        pyroute2 \
        python-iptables \
        git+https://github.com/larsks/python-netns.git \
        git+https://github.com/n0stack/n0library.git \
 && python setup.py install \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# development
RUN pip install --no-cache-dir \
        ipython
