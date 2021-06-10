#!/usr/local/bin/fish

# source the script with the below command
# source env.fish

echo "exporting environment variables for local development"
set -x MONGO_HOST mongo.default.svc.cluster.local
set -x MONGO_DATABASE wizard
set -x MONGO_PORT 27017
set -x GRPC_PORT 9999
set -x SERVER_HOST localhost
echo "done. happy coding!"