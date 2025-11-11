class DfCacheError(Exception):
    """Base exception for dfcache SDK."""


class NotFoundError(DfCacheError):
    """Raised when a cache entry does not exist (maps to os.ErrNotExist)."""


class DaemonUnavailableError(DfCacheError):
    """Raised when dfdaemon is not reachable or unhealthy."""
