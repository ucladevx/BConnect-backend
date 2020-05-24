CREATE ROLE connect_b WITH 
    LOGIN
    SUPERUSER
    INHERIT
    CREATEDB
    CREATEROLE
    REPLICATION
    PASSWORD 'connect_b';

GRANT ALL PRIVILEGES ON DATABASE connect_b_users TO connect_b;

CREATE DATABASE connect_b_test OWNER connect_b;

GRANT ALL PRIVILEGES ON DATABASE connect_b_test TO connect_b;