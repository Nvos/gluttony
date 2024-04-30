FROM postgres:16.2-bookworm

COPY ./docker/dictionary/polish.* /usr/share/postgresql/16/tsearch_data/
COPY ./docker/init.sh /docker-entrypoint-initdb.d