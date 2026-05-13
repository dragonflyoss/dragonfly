from __future__ import annotations

import os
import subprocess
import json
from typing import Optional

from .errors import DfCacheError, NotFoundError, DaemonUnavailableError

_D7Y_SCHEME = "d7y:/"  # prefix used to build internal URL from cid


def _cid_to_url(cid: str) -> str:
    from urllib.parse import quote
    return f"{_D7Y_SCHEME}{quote(cid, safe='') }"


def _run_cmd(args: list[str], timeout: float) -> subprocess.CompletedProcess:
    return subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, timeout=timeout)


class DfCacheClient:
    """Thin wrapper invoking existing dfcache CLI for stat/import/export/delete.

    This avoids needing Python gRPC stubs initially; later we can switch to gRPC.
    """

    def __init__(self, binary: Optional[str] = None, timeout: float = 10.0) -> None:
        self._binary = binary or self._detect_binary()
        self._timeout = timeout

    def _detect_binary(self) -> str:
        # Allow explicit override via env
        env_bin = os.getenv("DRAGONFLY_DFCACHE_BINARY")
        if env_bin and os.path.isfile(env_bin):
            return env_bin

        # Search common build output directories (linux/darwin, amd64/arm64)
        bin_dir = os.path.join(os.getcwd(), "bin")
        candidates: list[str] = []
        if os.path.isdir(bin_dir):
            for root, dirs, files in os.walk(bin_dir):
                if "dfcache" in files:
                    candidates.append(os.path.join(root, "dfcache"))
        # Fallback to PATH lookup
        candidates.append("dfcache")
        for c in candidates:
            if os.path.isfile(c) and os.access(c, os.X_OK):
                return c
        raise DaemonUnavailableError(
            "dfcache binary not found. Build it via 'make build-dfcache' or set DRAGONFLY_DFCACHE_BINARY."
        )

    def stat(self, cid: str, tag: str = "", local_only: bool = False, timeout: Optional[float] = None) -> bool:
        args = [self._binary, "stat", "-i", cid]
        if tag:
            args += ["-t", tag]
        if local_only:
            args += ["-l"]
        cp = _run_cmd(args, timeout or self._timeout)
        if cp.returncode == 0:
            return True
        if "not exist" in cp.stderr.lower() or "not exist" in cp.stdout.lower():
            return False
        if cp.returncode != 0:
            raise DfCacheError(f"stat failed: {cp.stderr.strip() or cp.stdout.strip()}")
        return False

    def import_cache(self, cid: str, path: str, tag: str = "", timeout: Optional[float] = None) -> None:
        if not os.path.isfile(path):
            raise FileNotFoundError(path)
        args = [self._binary, "import", "-i", cid, "-I", path]
        if tag:
            args += ["-t", tag]
        cp = _run_cmd(args, timeout or self._timeout)
        if cp.returncode != 0:
            raise DfCacheError(f"import failed: {cp.stderr.strip() or cp.stdout.strip()}")

    def export(self, cid: str, output: str, tag: str = "", local_only: bool = False, timeout: Optional[float] = None) -> None:
        parent = os.path.dirname(os.path.abspath(output))
        os.makedirs(parent, exist_ok=True)
        args = [self._binary, "export", "-i", cid, "-O", output]
        if tag:
            args += ["-t", tag]
        if local_only:
            args += ["-l"]
        cp = _run_cmd(args, timeout or self._timeout)
        if cp.returncode != 0:
            if "not exist" in cp.stderr.lower() or "not exist" in cp.stdout.lower():
                raise NotFoundError(f"cache {cid} not found")
            raise DfCacheError(f"export failed: {cp.stderr.strip() or cp.stdout.strip()}")

    def delete(self, cid: str, tag: str = "", timeout: Optional[float] = None) -> None:
        args = [self._binary, "delete", "-i", cid]
        if tag:
            args += ["-t", tag]
        cp = _run_cmd(args, timeout or self._timeout)
        if cp.returncode != 0 and "not exist" not in cp.stderr.lower():
            raise DfCacheError(f"delete failed: {cp.stderr.strip() or cp.stdout.strip()}")

    def check_health(self) -> bool:
        try:
            self.stat("health-check-cid")
            return True
        except DfCacheError:
            return False

    def info_json(self) -> str:
        data = {
            "binary": self._binary,
            "timeout": self._timeout,
        }
        return json.dumps(data)

