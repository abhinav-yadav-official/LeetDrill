#!/usr/bin/env python3
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PAGE = ROOT / "web" / "homepage" / "index.html"


def require(condition, message):
    if not condition:
        raise SystemExit(message)


def main():
    body = PAGE.read_text()
    require("https://github.com/abhinav-yadav-official/leetdrill" in body, "homepage must link GitHub repo")
    require("https://abhiy.xyz/leetdrill" in body, "homepage must link hosted LeetDrill")
    require("LeetDrill" in body, "homepage must mention LeetDrill")


if __name__ == "__main__":
    main()
