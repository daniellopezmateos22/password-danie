#!/bin/sh
set -e
: "${API_URL:=http://localhost:8080}"
echo "window.API_URL = '${API_URL}';" > /usr/share/nginx/html/config.js
exec nginx -g "daemon off;"
