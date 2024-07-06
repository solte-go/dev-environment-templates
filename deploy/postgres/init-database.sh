#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username postgres --dbname postgres <<-EOSQL
    CREATE USER um_operator PASSWORD 'dev_only_pwd';
    CREATE DATABASE user_management owner um_operator ENCODING = 'UTF-8';
EOSQL