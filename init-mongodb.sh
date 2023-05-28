#!/bin/bash
set -e

mongo <<EOF
use $MONGO_INITDB_DATABASE
requests = db.getSiblingDB('go-motivation-bot-db')
requests.createUser({
  user:  '$MONGO_USERNAME',
  pwd: '$MONGO_PASSWORD',
  roles: [{
    role: 'readWrite',
    db: 'go-motivation-bot-db'
  }]
})
EOF
