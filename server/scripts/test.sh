POSTGRES_USER="postgres"
POSTGRES_ADDR="localhost"
POSTGRES_DB_NAME="wedding_rsvp_test"

TESTCONTEXT=${1:--r}
printf "TESTING CONTEXT ${TESTCONTEXT}\n"

psql -h $POSTGRES_ADDR -U $POSTGRES_USER -c "DROP DATABASE IF EXISTS ${POSTGRES_DB_NAME};" && \
psql -h $POSTGRES_ADDR -U $POSTGRES_USER -c "CREATE DATABASE ${POSTGRES_DB_NAME};" && \
goose -env test up && \
WEDDING_RSVP_ENV=test \
ginkgo ${TESTCONTEXT}
