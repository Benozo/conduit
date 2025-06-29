#!/bin/bash
# Conduit MCP Server Launcher
# This script ensures the correct working directory when launching Conduit

cd "$(dirname "$0")"
exec go run main.go "$@"
