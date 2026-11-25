#!/usr/bin/env bash
# Copyright 2026 High Five Labs
# SPDX-License-Identifier: Apache-2.0
# Download the vendor's OpenAPI specs at the commits vendor.yaml pins.
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p specs
[ -s specs/cloudflare.json ] || curl -fsSL -o specs/cloudflare.json "https://raw.githubusercontent.com/cloudflare/api-schemas/8b1a451819bc546d410e41cca71a86bcc2779808/openapi.json"
