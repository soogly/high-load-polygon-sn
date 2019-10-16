#!/bin/bash

set -e

# install requirements
# go get github.com/lib/pq
# go get -u golang.org/x/crypto/bcrypt

# If you run that script not the first time 
su postgres -c "psql -d postgres < ./db/clear.sql"

# init db
su postgres -c "psql -U postgres < ./db/init.sql && \
                psql -U postgres -d go_app_db < ./db/tables.sql"

# build app
# go build main.go

# run app
./pain

exit 0
