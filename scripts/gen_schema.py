import requests
from pathlib import Path
import json
import re
from typing import List

root = Path(__file__).parent / '..'
schema_path = root / 'shortcut.schema.json'
with open(schema_path) as f:
    schema = json.load(f)

# finally add items types
table = requests.get('https://github.com/robotn/gohook/raw/c94ab299da47174a07a1d2d28d42750183062a51/tables.go').text
keys: List[str] = re.findall(r'\"(.+)\"', table)
keys = list(set(keys))
schema['properties']['shortcuts']['items']['properties']['keys']['items'] = {
    'type': 'string',
    'enum': keys
}
with open(schema_path, 'w') as f:
    json.dump(schema, f, indent=4)