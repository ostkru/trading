#!/bin/bash
cd /var/www/trading
export PATH="/root/.nvm/versions/node/v22.16.0/bin:$PATH"
exec npx redoc-cli serve openapi.json --port 8090 --host 0.0.0.0
