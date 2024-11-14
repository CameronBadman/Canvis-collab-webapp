kubectl create configmap cassandra-config --from-env-file=.env -n db


kubectl create secret generic cassandra-secrets \
  --from-literal=CASSANDRA_PASSWORD=cassandra1 \
  -n db

