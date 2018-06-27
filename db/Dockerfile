FROM mysql:8.0.11

ENV MYSQL_ROOT_PASSWORD 123456
ENV MYSQL_DATABASE blog

COPY ./up.sql /docker-entrypoint-initdb.d/

EXPOSE 3306
CMD ["mysqld"]
