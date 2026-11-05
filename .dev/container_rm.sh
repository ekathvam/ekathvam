
#!/usr/bin/env bash

set -Eeuo pipefail

CONTAINER="ekathvam"
LXC_PATH="/ekathvam/containers"

# Debian image cache used by Ekathvam.
CACHE_PATH="/var/cache/lxc/download/debian/trixie/amd64"

# Require root.
if [[ "$EUID" -ne 0 ]]; then
    echo "Error: Run with sudo."
    exit 1
fi

# Verify LXC installation.
if ! command -v lxc-info >/dev/null 2>&1; then
    echo "Error: LXC is not installed."
    exit 1
fi

echo "Cleaning Ekathvam development environment..."

# Remove the container if it exists.
if lxc-info -n "$CONTAINER" -P "$LXC_PATH" >/dev/null 2>&1; then

    STATE=$(lxc-info -n "$CONTAINER" -P "$LXC_PATH" -sH)

    if [[ "$STATE" == "RUNNING" ]]; then
        echo "Stopping container..."
        lxc-stop -n "$CONTAINER" -P "$LXC_PATH"
    fi

    echo "Destroying container..."
    lxc-destroy -n "$CONTAINER" -P "$LXC_PATH"

else
    echo "Container does not exist. Skipping."
fi

# Verify container removal.
if lxc-info -n "$CONTAINER" -P "$LXC_PATH" >/dev/null 2>&1; then
    echo "Error: Container still exists. Aborting cache cleanup."
    exit 1
fi

# Remove only the selected Debian image cache.
if [[ -d "$CACHE_PATH" ]]; then
    echo "Removing Debian Trixie AMD64 cache..."
    rm -rf -- "$CACHE_PATH"
else
    echo "Debian cache already clean."
fi

systemctl disable lxc-net

echo
echo "Ekathvam cleanup completed successfully."
echo "Container: removed"
echo "Debian image cache: removed if present"
echo "LXC installation: preserved"
echo "Project files: preserved"
echo "Storage parent: preserved"

