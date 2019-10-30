#!/bin/bash

set -e

# install requirements
# go get github.com/go-sql-driver/mysql
# go get -u golang.org/x/crypto/bcrypt

# If you run that script not the first time 
# su postgres -c "psql -d postgres < ./db/clear.sql"

# mysql < ./db/clear.sql

# init db
# su postgres -c "psql -U postgres < ./db/init.sql && \
#                 psql -U postgres -d go_app_db < ./db/tables.sql"

mysql < ./db/init.sql
mysql go_app_db < ./db/tables.sql
# build app
# go build main.go

# run app09(((((((((((((?KIO378793)))))))))))))
./pain

exit 0


