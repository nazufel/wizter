#!/bin/sh 

# source the script with the below command
# source env.sh

echo "exporting environment variables for local development"
export MONGO_HOST mongo.default.svc.cluster.local
export MONGO_DATABASE wizard
export MONGO_PORT 27017
export GRPC_PORT 9999
export SERVER_HOST localhost
echo "done. happy coding!"