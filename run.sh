#!/bin/bash

# Set repository owner and name
OWNER="vijitkamboj"
REPO="pprof-mcp"

# Fetch the latest release information
RELEASE_INFO=$(curl -s https://api.github.com/repos/$OWNER/$REPO/releases/latest)

# Extract release details
RELEASE_TAG=$(echo "$RELEASE_INFO" | grep -m 1 '"tag_name":' | cut -d '"' -f 4)

# Construct the URL for run.sh using the latest release tag
PROFILER_URL="https://github.com/$OWNER/$REPO/releases/download/$RELEASE_TAG/profiler"


echo "$PROFILER_URL"
# Download the latest profiler binary, make it executable and run it
curl -L "$PROFILER_URL" -o profiler
chmod +x profiler
./profiler