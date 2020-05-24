#!/bin/bash

DB=postgres
SQL_remover=removal

psql $DB < $SQL_remover.sql
