#!/bin/sh
set -e

# Resolve Supabase host IPs from host DNS and add to /etc/hosts
# This works around DNS issues on Windows with proxy/VPN

SUPABASE_POOLER="aws-1-eu-central-1.pooler.supabase.com"

# Try to resolve using host DNS via getent or nslookup
if command -v getent > /dev/null 2>&1; then
  IPS=$(getent hosts "$SUPABASE_POOLER" 2>/dev/null | awk '{print $1}')
elif command -v nslookup > /dev/null 2>&1; then
  IPS=$(nslookup "$SUPABASE_POOLER" 2>/dev/null | grep -oE '([0-9]{1,3}\.){3}[0-9]{1,3}' | grep -v "127\.\|8\.8\.8\|1\.1\.1" | tail -n +2)
fi

if [ -n "$IPS" ]; then
  for ip in $IPS; do
    if ! grep -q "$SUPABASE_POOLER" /etc/hosts; then
      echo "$ip $SUPABASE_POOLER" >> /etc/hosts
    fi
  done
  echo "entrypoint: added $SUPABASE_POOLER to /etc/hosts"
fi

exec "$@"
