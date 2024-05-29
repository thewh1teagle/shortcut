import fs from 'fs'
import path from 'path'

const root = path.join(__dirname, "..");
const schemaPath = path.join(root, "shortcut.schema.json");

// Read the schema file
const schema = JSON.parse(fs.readFileSync(schemaPath))

// Fetch the table data from the URL
const resp = await fetch("https://github.com/robotn/gohook/raw/c94ab299da47174a07a1d2d28d42750183062a51/tables.go")
if (!resp.ok) {
  throw new Error(`HTTP error! Status: ${resp.status}`);
}
const table = await resp.text()
// Extract keys using regex
const keys = Array.from(
  new Set(table.match(/"(.+?)"/g).map((match) => match.slice(1, -1)))
);
// Sort alphabet
keys.sort((a, b) => a.localeCompare(b));
console.log(keys);
// Update schema
schema.properties.shortcuts.items.properties.keys.items = {
  type: "string",
  enum: keys,
};
// Write updated schema back to file
fs.writeFileSync(schemaPath, JSON.stringify(schema, null, 4));
console.log("Schema updated successfully.");