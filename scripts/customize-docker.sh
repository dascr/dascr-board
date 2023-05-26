#!/bin/bash
echo "[*] Checking customization through frontend/.env"
if [ "$VITE_API_BASE" = "http://localhost:8000/" ]; then
    echo "[*] VITE_API_BASE setting is default (localhost)"
else
	echo "[*] VITE_API_BASE custom setting - $VITE_API_BASE"
	sed -i -r "s#http://localhost:8000/\$#$VITE_API_BASE#" frontend/Dockerfile
fi

if [ "$VITE_API_URL" = "http://localhost:8000/api" ]; then
    echo "[*] VITE_API_URL setting is default (localhost)"
else
	echo "[*] VITE_API_URL custom setting - $VITE_API_URL"
	sed -i -r "s#http://localhost:8000/api\$#$VITE_API_URL#" frontend/Dockerfile
fi

if [ "$VITE_WS_URL" = "ws://localhost:8000/ws" ]; then
    echo "[*] VITE_WS_URL setting is default (localhost)"
else
	echo "[*] VITE_WS_URL custom setting - $VITE_WS_URL"
	sed -i -r "s#ws://localhost:8000/ws\$#$VITE_WS_URL#" frontend/Dockerfile
fi
