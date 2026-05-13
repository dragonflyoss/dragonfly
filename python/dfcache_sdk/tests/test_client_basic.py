import os
import tempfile
import shutil
import stat
import subprocess
import sys
from pathlib import Path

import pytest

from dragonfly_dfcache import DfCacheClient, NotFoundError


def make_fake_dfcache(tmp: Path) -> str:
    fake = tmp / "dfcache"
    fake.write_text(
        """#!/usr/bin/env bash
cmd=$1; shift
case "$cmd" in
  stat)
    # simulate not-exist if CID contains 'missing'
    if [[ "$@" == *missing* ]]; then
      echo "not exist" >&2
      exit 1
    fi
    exit 0
    ;;
  import)
    exit 0
    ;;
  export)
    # write output file path passed after -O
    while [[ $# -gt 0 ]]; do
      if [[ "$1" == "-O" ]]; then
        shift; out=$1; break
      fi
      shift
    done
    mkdir -p "$(dirname "$out")" && echo ok > "$out"
    exit 0
    ;;
  delete)
    exit 0
    ;;
  *)
    echo "unknown" >&2
    exit 2
    ;;
 esac
"""
    )
    st = os.stat(fake)
    os.chmod(fake, st.st_mode | stat.S_IEXEC)
    return str(fake)


def test_basic_flow(tmp_path: Path, monkeypatch: pytest.MonkeyPatch):
    fake = make_fake_dfcache(tmp_path)
    monkeypatch.setenv("DRAGONFLY_DFCACHE_BINARY", fake)

    c = DfCacheClient(timeout=2)

    # stat existing
    assert c.stat("present-cid") is True
    # stat missing
    assert c.stat("missing-cid") is False

    # import requires file
    src = tmp_path / "file.bin"
    src.write_bytes(b"x")
    c.import_cache("present-cid", str(src))

    # export should create output
    out = tmp_path / "out.bin"
    c.export("present-cid", str(out))
    assert out.exists()

    # delete
    c.delete("present-cid")
