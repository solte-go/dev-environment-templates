#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username postgres --dbname postgres <<-EOSQL
    CREATE USER operator PASSWORD 'dev_only_pwd';
    CREATE DATABASE local_dev owner operator ENCODING = 'UTF-8';
EOSQL