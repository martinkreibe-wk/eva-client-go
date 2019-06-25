#!/bin/bash
# Updates table of contents in all relevant readmes
# Requires markdown-toc
MAX_DEPTH=4

# Parent README
npx markdown-toc -i README.md --maxdepth $MAX_DEPTH
npx markdown-toc -i edn/README.md --maxdepth $MAX_DEPTH
npx markdown-toc -i eva/http/README.md --maxdepth $MAX_DEPTH
