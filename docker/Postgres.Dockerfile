FROM postgres:16

ADD /scripts/postgres/create_tables.sql /docker-entrypoint-initdb.d
ADD /scripts/postgres/seed_data.sql /docker-entrypoint-initdb.d

RUN chmod a+r /docker-entrypoint-initdb.d/*