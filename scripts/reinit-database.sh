#!/bin/bash

PSQLCMD="psql postgresql://root@localhost:5432/wallawire?sslmode=verify-full&sslcert=$PWD/walladata/certs/dbclient/client.root.crt&sslkey=$PWD/walladata/certs/dbclient/client.root.key&sslrootcert=$PWD/walladata/certs/dbclient/ca.crt"

${PSQLCMD} < schema/02_down.sql
${PSQLCMD} < schema/01_down.sql
${PSQLCMD} < schema/01_up.sql
${PSQLCMD} < schema/02_up.sql
