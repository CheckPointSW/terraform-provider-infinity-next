#!/usr/bin/env bash
TOKEN=$(terraform output access_token | jq -r .)

ENFORCE_RESPONSE=$(curl --location --request POST 'https://cloudinfra-gw.portal.checkpoint.com/app/i2/graphql/V1' --header 'Authorization: Bearer '$TOKEN'' --header 'Content-Type: application/json' --data-raw '{"query":"mutation {\r\n    enforcePolicy {\r\n        id\r\n    }\r\n}","variables":{}}')

if [ "$(jq 'has(".data.enforcePolicy.id")' <<< $ENFORCE_RESPONSE)" == "true" ]; then
printf "Successfuly enforced policy\nTask ID=$(jq -r .data.enforcePolicy.id <<< $ENFORCE_RESPONSE)\n"
else
printf "Failed to enforce policy: $(jq -r .message <<< $ENFORCE_RESPONSE)\n"
fi
