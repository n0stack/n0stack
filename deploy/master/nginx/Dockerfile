FROM swaggerapi/swagger-ui AS SWAGGER

ENV API_URL=n0stack.swagger.json
COPY swagger_run.sh /swagger_run.sh

RUN sh /swagger_run.sh

FROM nginx

COPY --from=SWAGGER /usr/share/nginx/html /usr/share/nginx/html/api
RUN cd /usr/share/nginx/html/api \
 && apt update \
 && apt install -y wget \
 && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* \
 && wget https://raw.githubusercontent.com/n0stack/n0stack/master/n0proto.swagger.json/n0stack.swagger.json

COPY etc/nginx/conf.d /etc/nginx/conf.d
