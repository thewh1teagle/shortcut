const fs = require("fs");
const path = require("path");

const root = path.join(__dirname, "..");
const schemaPath = path.join(root, "shortcut.schema.json");

// Read the schema file
const schema = require(schemaPath);

// Fetch the table data from the URL
fetch(
  "https://github.com/robotn/gohook/raw/c94ab299da47174a07a1d2d28d42750183062a51/tables.go"
)
  .then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.text();
  })
  .then((table) => {
    // Extract keys using regex
    const keys = Array.from(
      new Set(table.match(/"(.+?)"/g).map((match) => match.slice(1, -1)))
    );
    console.log(keys);
    // Update schema
    schema.properties.shortcuts.items.properties.keys.items = {
      type: "string",
      enum: keys,
    };

    // Write updated schema back to file
    fs.writeFileSync(schemaPath, JSON.stringify(schema, null, 4));

    console.log("Schema updated successfully.");
  })
  .catch((error) => {
    console.error("Error fetching data:", error);
  });
