#!/bin/bash

PSQLCMD="psql postgresql://root@localhost:5432/wallawire?sslmode=verify-full&sslcert=$PWD/walladata/certs/dbclient/client.root.crt&sslkey=$PWD/walladata/certs/dbclient/client.root.key&sslrootcert=$PWD/walladata/certs/dbclient/ca.crt"

${PSQLCMD} <<EOF
CREATE USER wallawire;
CREATE DATABASE wallawire;
GRANT SELECT,INSERT,UPDATE,DELETE ON DATABASE wallawire TO wallawire;
EOF

${PSQLCMD} < schema/01_up.sql
${PSQLCMD} < schema/02_up.sql
