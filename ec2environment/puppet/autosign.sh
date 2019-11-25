#!/bin/bash

CSR=$(< /dev/stdin)
CERTNAME=$1

function extension {
  echo "$csr" | openssl req -noout -text | fgrep -A1 "$1" | tail -n 1 \
      | sed -e 's/^ *//;s/ *$//'
}

psk=$(extension '1.3.6.1.4.1.34380.1.1.4')

if [ -f "/etc/puppetlabs/puppet/psk.txt" ]; then
    if grep -q "$psk" "/etc/puppetlabs/puppet/psk.txt"; then
        exit 0
    else
        exit 1
    fi
else
    echo "Could not find PSK file for $CERTNAME"
    exit 1
fi

exit 1