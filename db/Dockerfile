FROM postgres:10.3

ENV POSTGRES_DB blog
ENV POSTGRES_USER blog
ENV POSTGRES_PASSWORD 123456

COPY ./up.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]
