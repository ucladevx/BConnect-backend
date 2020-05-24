#!/bin/bash

DB=postgres
SQL_remover=removal

psql $DB < $GOPATH/src/github.com/ucladevx/BConnect-backend/scripts/$SQL_remover.sql
