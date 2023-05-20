#!/bin/bash
set -e

mongo <<EOF
use $MONGO_INITDB_DATABASE
requests = db.getSiblingDB('go-motivation-bot-users')
requests.createUser({
  user:  '$MONGO_USERNAME',
  pwd: '$MONGO_PASSWORD',
  roles: [{
    role: 'readWrite',
    db: 'go-motivation-bot-users'
  }]
})
EOF
