#!/bin/bash

DB=postgres
SQL_intiliazer=setup
SQL_remover=removal

if psql $DB -t -c '\du' | cut -d \| -f 1 | grep -qw 'connect_b'; then
    echo INITIALIZATION ALREADY COMPLETE...
else
    if psql -lqt | cut -d \| -f 1 | grep -qw 'connect_b_users'; then
        psql $DB < $SQL_remover.sql
        heroku pg:pull DATABASE connect_b_users --app protected-refuge-33249
        psql $DB < $SQL_intiliazer.sql
        echo DATABASE INITIALIZED...
    else
        heroku pg:pull DATABASE connect_b_users --app protected-refuge-33249
        psql $DB < $SQL_intiliazer.sql
        echo DATABASE INITIALIZED...
    fi
fi