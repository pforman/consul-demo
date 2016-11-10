FROM ubuntu:latest

RUN apt-get update && apt-get install -y nginx supervisor
RUN mkdir -p /var/log/supervisor

COPY consul-cli.linux /usr/bin/consul-cli
COPY consul-lock.bash /usr/bin
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 80

CMD ["/usr/bin/supervisord"]
