FROM python:3.6

COPY . /srv
WORKDIR /srv
RUN apt-get update \
 && apt-get install -y \
        dnsmasq \
        iptables \
 && pip install --no-cache-dir -r requirements.txt \
 && python setup.py install \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# development
RUN pip install --no-cache-dir \
        ipython
