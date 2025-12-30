#!/usr/bin/env node

/**
 * Checks if the current package.json version is newer than what's in the API.
 * Exits with code 0 if should release, 1 if not.
 */

const pkg = require('../package.json');
const https = require('https');

const API_URL = process.env.API_URL;

if (!API_URL) {
  console.error('API_URL environment variable is required');
  process.exit(1);
}

function compareVersions(a, b) {
  const partsA = a.replace(/^v/, '').split('.').map(Number);
  const partsB = b.replace(/^v/, '').split('.').map(Number);

  for (let i = 0; i < Math.max(partsA.length, partsB.length); i++) {
    const numA = partsA[i] || 0;
    const numB = partsB[i] || 0;
    if (numA > numB) return 1;
    if (numA < numB) return -1;
  }
  return 0;
}

async function getLatestVersion() {
  return new Promise((resolve) => {
    const url = `${API_URL}/desktop/releases/latest`;
    
    https.get(url, {
      headers: {
        'User-Agent': 'Aniways-Desktop-Release/1.0',
      },
    }, (res) => {
      if (res.statusCode === 404) {
        // No releases yet
        resolve(null);
        return;
      }
      
      if (res.statusCode !== 200) {
        console.error(`API returned status ${res.statusCode}`);
        resolve(null);
        return;
      }

      let data = '';
      res.on('data', (chunk) => (data += chunk));
      res.on('end', () => {
        try {
          const json = JSON.parse(data);
          resolve(json.version);
        } catch {
          resolve(null);
        }
      });
    }).on('error', (err) => {
      console.error('Failed to fetch latest version:', err.message);
      resolve(null);
    });
  });
}

async function main() {
  const currentVersion = pkg.version;
  console.log(`Current package.json version: ${currentVersion}`);

  const latestVersion = await getLatestVersion();
  
  if (latestVersion === null) {
    console.log('No existing releases found. Should release.');
    console.log(`::set-output name=should_release::true`);
    console.log(`::set-output name=version::${currentVersion}`);
    process.exit(0);
  }

  console.log(`Latest released version: ${latestVersion}`);

  if (compareVersions(currentVersion, latestVersion) > 0) {
    console.log(`Version ${currentVersion} is newer than ${latestVersion}. Should release.`);
    console.log(`::set-output name=should_release::true`);
    console.log(`::set-output name=version::${currentVersion}`);
    process.exit(0);
  } else {
    console.log(`Version ${currentVersion} is not newer than ${latestVersion}. Skipping release.`);
    console.log(`::set-output name=should_release::false`);
    process.exit(0);
  }
}

main();
