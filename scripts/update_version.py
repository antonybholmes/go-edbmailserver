import json
from dotenv import dotenv_values
from datetime import datetime

# current time and date
# datetime object
time = datetime.now()

config = json.load(open("version.json"))

# print(config)

major, minor, patch, build = [int(x) for x in config["version"].split(".")]

build += 1
patch += 1

if patch > 9:
    minor += 1
    patch = 0

if minor > 9:
    major += 1
    minor = 0

config["version"] = f"{major}.{minor}.{patch}.{build}"
config["updated"] = time.strftime("%b %d, %Y")

with open("version.json", "w") as f:
    json.dump(config, f, indent=4)

# with open("version.env", "w") as f:
#     for key in sorted(config):
#         print(f'{key}="{config[key]}"', file=f)
