#!/usr/bin/env python3

import json
import re
import sys

payload = json.load(sys.stdin)

if payload.get("tool_name") != "Bash":
    sys.exit(0)

command = payload.get("tool_input", {}).get("command", "").strip()

RULES = [
    (
        re.compile(r"^find\s+\.(?:\s|$)"),
        "Avoid repository-wide searches. Search a specific subtree instead (e.g. `find src ...` or `find packages/api ...`).",
    ),
    (
        re.compile(r"^rg\s+--files\s*$"),
        "Avoid listing the entire repository. Use `rg --files <subtree>/` instead (e.g. `src/` or `packages/api/`).",
    ),
    (
        re.compile(r"^fd(?:\s+\.)?\s*$"),
        "Avoid enumerating the entire repository. Scope `fd` to a specific directory.",
    ),
    (
        re.compile(r"^tree(?:\s+\.)?\s*$"),
        "Avoid printing the entire repository tree. Inspect only the relevant subtree.",
    ),
    (
        re.compile(r"^ls\s+-R(?:\s+\.)?\s*$"),
        "Avoid recursive repository listings. Inspect only the relevant subtree.",
    ),
]

for pattern, reason in RULES:
    if pattern.fullmatch(command) or pattern.match(command):
        print(json.dumps({
            "hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "deny",
                "permissionDecisionReason": reason,
            }
        }))
        sys.exit(0)

# Allow all other commands.
sys.exit(0)
