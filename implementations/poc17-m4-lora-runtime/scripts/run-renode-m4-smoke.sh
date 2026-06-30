#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
FIRMWARE_DIR="$ROOT_DIR/firmware/renode-m4-smoke"

docker run --rm \
	-v "$FIRMWARE_DIR:/work" \
	-w /work \
	rust:1.84 \
	sh -lc 'export PATH="/usr/local/cargo/bin:$PATH"; rustup target add thumbv7em-none-eabihf && cargo build --release'

docker run --rm \
	-v "$ROOT_DIR:/work" \
	-w /work \
	antmicro/renode:latest \
	renode --disable-gui --console renode/poc17-m4-smoke.resc
