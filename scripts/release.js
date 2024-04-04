const fs = require("fs");
const path = require("path");
const child_process = require("child_process");

// Github options
const OWNER = "thewh1teagle";
const REPO = "shortcut";

// Paths
const ROOT = path.join(__dirname, "..");
const MAIN = path.join(ROOT, "main.go");

// Environment
const OS = {
  win32: "windows",
  linux: "linux",
  darwin: "macos",
}[process.platform];
const ARCH = process.arch;
const VERSION = fs
  .readFileSync(path.join(ROOT, "main.go"), "utf-8")
  .match(/version\s*=\s*"([^"]+)"/)[1];
const TOKEN = process.env.GITHUB_TOKEN;

// Dist
const NAME = `shortcut-${OS}-${ARCH}-${VERSION}`;
const DIST = NAME + (OS == "windows" ? ".zip" : ".tar.gz");

// Prepare
fs.rmSync(path.join(ROOT, NAME), { recursive: true, force: true });
fs.mkdirSync(path.join(ROOT, NAME));

// Build
child_process.execSync("go build -tags release", { cwd: ROOT, shell: true });

// Bundle
console.log(fs.readdirSync(ROOT))
fs.renameSync(path.join(ROOT, "shortcut"), path.join(ROOT, NAME, 'shortcut'));
fs.copyFileSync(
  path.join(ROOT, "shortcut.conf.json"),
  path.join(ROOT, NAME, `shortcut.conf.json`)
);
fs.copyFileSync(
  path.join(ROOT, "shortcut.schema.json"),
  path.join(ROOT, NAME, 'shortcut.schema.json')
);

// Compress
if (OS == "windows") {
  child_process.execSync(`7z a ${DIST} ${NAME}`, { cwd: ROOT });
} else {
  child_process.execSync(`tar -czvf ${DIST} ${NAME}`, { cwd: ROOT });
}

// Clean
fs.rmSync(NAME, { recursive: true });

// Publish

publish()

async function checkResponse(response) {
    if (![200,201,204].includes(response.status)) {
        const reason = await response.text()
        console.error(`status ${response.status} for ${response.url}: `, reason)
        process.exit(1)    
    }
}

async function publish() {
  try {
    const res = await fetch(`https://api.github.com/repos/${OWNER}/${REPO}/releases`, {
      headers: {
        Authorization: `Bearer ${TOKEN}`,
        "X-GitHub-Api-Version": "2022-11-28",
        Accept: "application/vnd.github+json",
      },
      method: "POST",
      body: JSON.stringify({
        tag_name: `v${VERSION}`,
        target_commitish: "main",
        name: `v${VERSION}`,
        body: "See assets for download",
        draft: false,
        prerelease: false,
        generate_release_notes: false,
      }),
    });
    const data = await res.json()
    if (!data?.errors?.code === 'already_exists') {
        await checkResponse(res)
    }
    
    console.log(`Created Release ${VERSION}`);
  } catch (e) {
    console.error(`Failed to create release ${VERSION}: ${e}`);
  }

  // Create Release

  // Delete previous asset if exists
  // curl -L \
  //   -H "Accept: application/vnd.github+json" \
  //   -H "Authorization: Bearer <YOUR-TOKEN>" \
  //   -H "X-GitHub-Api-Version: 2022-11-28" \
  //   https://api.github.com/repos/OWNER/REPO/releases/RELEASE_ID/assets
  // curl -L \
  //   -X DELETE \
  //   -H "Accept: application/vnd.github+json" \
  //   -H "Authorization: Bearer <YOUR-TOKEN>" \
  //   -H "X-GitHub-Api-Version: 2022-11-28" \
  //   https://api.github.com/repos/OWNER/REPO/releases/assets/ASSET_ID

  // Get release ID
  const res = await fetch(`https://api.github.com/repos/${OWNER}/${REPO}/releases/tags/v${VERSION}`, {
    headers: {
        Authorization: `Bearer ${TOKEN}`,
        "X-GitHub-Api-Version": "2022-11-28",
        Accept: "application/vnd.github+json",
    }
  })
  const releaseData = await res.json()
  const releaseID = releaseData.id
  const prevID = releaseData.assets.find(a => a.name === DIST)?.id

  if (prevID) {
    // Delete previous asset
    const res = await fetch(`https://api.github.com/repos/${OWNER}/${REPO}/releases/assets/${prevID}`, {
      headers: {
        Authorization: `Bearer ${TOKEN}`,
        "X-GitHub-Api-Version": "2022-11-28",
        Accept: "application/vnd.github+json",
      },
      method: "DELETE",
    });
    await checkResponse(res)
  }

  // Upload asset
  try {
    const res = await fetch(
      `https://uploads.github.com/repos/${OWNER}/${REPO}/releases/${releaseID}/assets?name=${DIST}`,
      {
        headers: {
          Authorization: `Bearer ${TOKEN}`,
          "X-GitHub-Api-Version": "2022-11-28",
          Accept: "application/vnd.github+json",
          "Content-Type": "application/octet-stream",
        },
        method: "POST",
        body: fs.readFileSync(path.join(ROOT, DIST)),
      }
    );
    await checkResponse(res)
    console.log(`Upload asset ${DIST} successfuly!`);
  } catch (e) {
    console.error(`Failed to upload asset ${DIST}: ${e}`);
  }
}
