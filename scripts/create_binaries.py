import json
import requests
import re
import sys
import logging

logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")

TAG_PATTERN = r"^v\d+\.\d+\.\d+$"
SUPPORTED_PLATFORMS = ["linux-amd64", "linux-arm64", "darwin-amd64", "darwin-arm64"]

def validate_tag(tag: str) -> bool:
    """Validate the tag format as vX.Y.Z."""
    return bool(re.match(TAG_PATTERN, tag))

def fetch_checksums(tag: str) -> dict:
    """Fetch and parse the sha256sum.txt file for the given tag."""
    if not validate_tag(tag):
        logging.error("The provided tag '%s' does not follow the 'vX.Y.Z' format.", tag)
        sys.exit(1)

    url = f"https://github.com/GGEZLabs/ggezchain/releases/download/{tag}/sha256sum.txt"
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
    except requests.RequestException as e:
        logging.error("Error fetching sha256sum.txt: %s", e)
        sys.exit(1)

    checksums = {}
    for line in response.text.splitlines():
        match = re.match(r"(\w+)\s+ggezchaind-[^\s]+-([\w]+-[\w]+)\.tar\.gz$", line)
        if match:
            checksum, platform = match.groups()
            if platform in SUPPORTED_PLATFORMS:
                checksums[platform] = checksum

    if not checksums:
        logging.warning("No valid checksums found in sha256sum.txt for supported platforms.")

    return checksums

def generate_info(checksums: dict, tag: str) -> dict:
    """Generate the info dictionary with binary URLs and checksums."""
    base_url = f"https://github.com/GGEZLabs/ggezchain/releases/download/{tag}/ggezchaind-{tag}-{{}}.tar.gz?checksum=sha256:{{}}"

    binaries = {}
    for platform in SUPPORTED_PLATFORMS:
        if platform in checksums:
            binaries[platform.replace("-", "/")] = base_url.format(platform, checksums[platform])
        else:
            logging.warning("No checksum found for platform: %s", platform)

    if not binaries:
        logging.error("No binaries were generated. Please check the sha256sum file.")
        sys.exit(1)

    return {"binaries": binaries}

def format_binaries(tag: str):
    """Fetch checksums, build binary info, and print as escaped JSON string."""
    try:
        checksums = fetch_checksums(tag)
        info = generate_info(checksums, tag)
        print(json.dumps(json.dumps(info)))
    except Exception as e:
        logging.error("Unexpected error updating JSON: %s", e)
        sys.exit(1)

def main():
    if len(sys.argv) != 2:
        print("Usage: python update_json_info.py <tag>")
        sys.exit(1)

    tag = sys.argv[1]
    format_binaries(tag)

if __name__ == "__main__":
    main()
