#!/usr/bin/env node

/**
 * Checks if the current package.json version is newer than what's in the API.
 * Writes outputs to GITHUB_OUTPUT file.
 */

const pkg = require('../package.json');
const https = require('https');
const http = require('http');
const fs = require('fs');

const API_URL = process.env.API_URL;
const GITHUB_OUTPUT = process.env.GITHUB_OUTPUT;

if (!API_URL) {
  console.error('API_URL environment variable is required');
  process.exit(1);
}

function writeOutput(name, value) {
  console.log(`${name}=${value}`);
  if (GITHUB_OUTPUT) {
    fs.appendFileSync(GITHUB_OUTPUT, `${name}=${value}\n`);
  }
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
    const client = url.startsWith('https') ? https : http;
    
    const req = client.get(url, {
      headers: {
        'User-Agent': 'Aniways-Desktop-Release/1.0',
      },
      timeout: 30000,
    }, (res) => {
      if (res.statusCode === 404) {
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
    });

    req.on('error', (err) => {
      console.error('Failed to fetch latest version:', err.message);
      resolve(null);
    });

    req.on('timeout', () => {
      console.error('Request timed out');
      req.destroy();
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
    writeOutput('should_release', 'true');
    writeOutput('version', currentVersion);
    process.exit(0);
  }

  console.log(`Latest released version: ${latestVersion}`);

  if (compareVersions(currentVersion, latestVersion) > 0) {
    console.log(`Version ${currentVersion} is newer than ${latestVersion}. Should release.`);
    writeOutput('should_release', 'true');
    writeOutput('version', currentVersion);
    process.exit(0);
  } else {
    console.log(`Version ${currentVersion} is not newer than ${latestVersion}. Skipping release.`);
    writeOutput('should_release', 'false');
    process.exit(0);
  }
}

main();
