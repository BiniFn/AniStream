#!/usr/bin/env node

/**
 * Uploads built artifacts to R2 and registers them in the API.
 * 
 * Required environment variables:
 * - R2_ACCOUNT_ID: Cloudflare account ID
 * - R2_ACCESS_KEY_ID: R2 access key
 * - R2_SECRET_ACCESS_KEY: R2 secret key
 * - R2_BUCKET_NAME: R2 bucket name
 * - R2_PUBLIC_URL: Public URL for the R2 bucket (e.g., https://releases.aniways.xyz)
 * - API_URL: API base URL
 * - DESKTOP_RELEASE_KEY: API key for creating releases
 */

const fs = require('fs');
const path = require('path');
const https = require('https');
const crypto = require('crypto');

const pkg = require('../package.json');

// Required env vars
const requiredEnvVars = [
  'R2_ACCOUNT_ID',
  'R2_ACCESS_KEY_ID', 
  'R2_SECRET_ACCESS_KEY',
  'R2_BUCKET_NAME',
  'R2_PUBLIC_URL',
  'API_URL',
  'DESKTOP_RELEASE_KEY',
];

for (const envVar of requiredEnvVars) {
  if (!process.env[envVar]) {
    console.error(`Missing required environment variable: ${envVar}`);
    process.exit(1);
  }
}

const {
  R2_ACCOUNT_ID,
  R2_ACCESS_KEY_ID,
  R2_SECRET_ACCESS_KEY,
  R2_BUCKET_NAME,
  R2_PUBLIC_URL,
  API_URL,
  DESKTOP_RELEASE_KEY,
} = process.env;

const VERSION = pkg.version;
const DIST_DIR = path.join(__dirname, '..', 'release');

// Platform mapping: file pattern -> platform key
const PLATFORM_PATTERNS = [
  { pattern: /-mac-arm64\.zip$/, platform: 'darwin-arm64', type: 'zip' },
  { pattern: /-mac-x64\.zip$/, platform: 'darwin-x64', type: 'zip' },
  { pattern: /-arm64\.exe$/, platform: 'win32-arm64', type: 'exe' },
  { pattern: / Setup .*\.exe$/, platform: 'win32-x64', type: 'exe' },
  { pattern: /-arm64\.AppImage$/, platform: 'linux-arm64', type: 'AppImage' },
  { pattern: /-x86_64\.AppImage$/, platform: 'linux-x64', type: 'AppImage' },
  { pattern: /\.AppImage$/, platform: 'linux-x64', type: 'AppImage' },
];

// Also upload these files for electron-updater
const UPDATER_FILES = [
  'latest.yml',
  'latest-mac.yml', 
  'latest-linux.yml',
];

function sign(key, msg) {
  return crypto.createHmac('sha256', key).update(msg).digest();
}

function getSignatureKey(key, dateStamp, regionName, serviceName) {
  const kDate = sign('AWS4' + key, dateStamp);
  const kRegion = sign(kDate, regionName);
  const kService = sign(kRegion, serviceName);
  const kSigning = sign(kService, 'aws4_request');
  return kSigning;
}

async function uploadToR2(filePath, key, originalFileName) {
  const fileContent = fs.readFileSync(filePath);
  const fileSize = fs.statSync(filePath).size;
  
  const host = `${R2_ACCOUNT_ID}.r2.cloudflarestorage.com`;
  const region = 'auto';
  const service = 's3';
  const method = 'PUT';
  
  // Content-Disposition header for proper download filename
  const contentDisposition = `attachment; filename="${originalFileName}"`;
  
  const now = new Date();
  const amzDate = now.toISOString().replace(/[:-]|\.\d{3}/g, '');
  const dateStamp = amzDate.slice(0, 8);
  
  const contentHash = crypto.createHash('sha256').update(fileContent).digest('hex');
  
  // URI-encode the key for the path, but use unencoded for canonical URI signing
  const canonicalUri = `/${R2_BUCKET_NAME}/${key.split('/').map(encodeURIComponent).join('/')}`;
  const encodedPath = `/${R2_BUCKET_NAME}/${key.split('/').map(encodeURIComponent).join('/')}`;
  
  const canonicalQuerystring = '';
  const canonicalHeaders = `content-disposition:${contentDisposition}\nhost:${host}\nx-amz-content-sha256:${contentHash}\nx-amz-date:${amzDate}\n`;
  const signedHeaders = 'content-disposition;host;x-amz-content-sha256;x-amz-date';
  
  const canonicalRequest = `${method}\n${canonicalUri}\n${canonicalQuerystring}\n${canonicalHeaders}\n${signedHeaders}\n${contentHash}`;
  
  const algorithm = 'AWS4-HMAC-SHA256';
  const credentialScope = `${dateStamp}/${region}/${service}/aws4_request`;
  const stringToSign = `${algorithm}\n${amzDate}\n${credentialScope}\n${crypto.createHash('sha256').update(canonicalRequest).digest('hex')}`;
  
  const signingKey = getSignatureKey(R2_SECRET_ACCESS_KEY, dateStamp, region, service);
  const signature = crypto.createHmac('sha256', signingKey).update(stringToSign).digest('hex');
  
  const authorizationHeader = `${algorithm} Credential=${R2_ACCESS_KEY_ID}/${credentialScope}, SignedHeaders=${signedHeaders}, Signature=${signature}`;
  
  return new Promise((resolve, reject) => {
    const req = https.request({
      hostname: host,
      path: encodedPath,
      method: 'PUT',
      headers: {
        'Host': host,
        'Content-Disposition': contentDisposition,
        'x-amz-date': amzDate,
        'x-amz-content-sha256': contentHash,
        'Authorization': authorizationHeader,
        'Content-Length': fileSize,
      },
    }, (res) => {
      let data = '';
      res.on('data', (chunk) => (data += chunk));
      res.on('end', () => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          // Return URL with encoded key
          resolve(`${R2_PUBLIC_URL}/${key.split('/').map(encodeURIComponent).join('/')}`);
        } else {
          reject(new Error(`Upload failed with status ${res.statusCode}: ${data}`));
        }
      });
    });
    
    req.on('error', reject);
    req.write(fileContent);
    req.end();
  });
}

async function registerRelease(platform, downloadUrl, fileName, fileSize) {
  const body = JSON.stringify({
    version: VERSION,
    platform,
    downloadUrl,
    fileName,
    fileSize,
    releaseNotes: '',
  });

  return new Promise((resolve, reject) => {
    const url = new URL(`${API_URL}/desktop/releases`);
    
    const req = https.request({
      hostname: url.hostname,
      port: url.port || 443,
      path: url.pathname,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${DESKTOP_RELEASE_KEY}`,
        'Content-Length': Buffer.byteLength(body),
        'User-Agent': 'Aniways-Desktop-Release/1.0',
      },
    }, (res) => {
      let data = '';
      res.on('data', (chunk) => (data += chunk));
      res.on('end', () => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(JSON.parse(data));
        } else {
          reject(new Error(`API call failed with status ${res.statusCode}: ${data}`));
        }
      });
    });

    req.on('error', reject);
    req.write(body);
    req.end();
  });
}

function detectPlatform(fileName) {
  for (const { pattern, platform } of PLATFORM_PATTERNS) {
    if (pattern.test(fileName)) {
      return platform;
    }
  }
  return null;
}

async function main() {
  console.log(`Uploading release v${VERSION}`);

  if (!fs.existsSync(DIST_DIR)) {
    console.error(`Release directory not found: ${DIST_DIR}`);
    process.exit(1);
  }

  const files = fs.readdirSync(DIST_DIR);
  const uploadedPlatforms = new Set();

  // Upload updater manifest files first
  for (const updaterFile of UPDATER_FILES) {
    const filePath = path.join(DIST_DIR, updaterFile);
    if (fs.existsSync(filePath)) {
      console.log(`Uploading updater file: ${updaterFile}`);
      try {
        const key = updaterFile;
        await uploadToR2(filePath, key, updaterFile);
        console.log(`  ✓ Uploaded ${updaterFile} to ${R2_PUBLIC_URL}/${updaterFile}`);
      } catch (err) {
        console.error(`  ✗ Failed to upload ${updaterFile}:`, err.message);
      }
    } else {
      console.warn(`  ⚠ Updater manifest not found: ${updaterFile} (this is normal for some platforms)`);
    }
  }

  // Upload release binaries
  for (const fileName of files) {
    const platform = detectPlatform(fileName);
    if (!platform) continue;

    // Skip if we already uploaded for this platform (prefer first match)
    if (uploadedPlatforms.has(platform)) continue;

    const filePath = path.join(DIST_DIR, fileName);
    const fileSize = fs.statSync(filePath).size;

    console.log(`\nProcessing: ${fileName}`);
    console.log(`  Platform: ${platform}`);
    console.log(`  Size: ${(fileSize / 1024 / 1024).toFixed(2)} MB`);

    try {
      // Upload to R2
      const key = `${VERSION}/${fileName}`;
      console.log(`  Uploading to R2...`);
      const downloadUrl = await uploadToR2(filePath, key, fileName);
      console.log(`  ✓ Uploaded to ${downloadUrl}`);

      // Register in API
      console.log(`  Registering in API...`);
      await registerRelease(platform, downloadUrl, fileName, fileSize);
      console.log(`  ✓ Registered in API`);

      uploadedPlatforms.add(platform);
    } catch (err) {
      console.error(`  ✗ Failed:`, err.message);
      process.exit(1);
    }
  }

  console.log(`\n✓ Successfully uploaded ${uploadedPlatforms.size} platform(s)`);
}

main().catch((err) => {
  console.error('Fatal error:', err);
  process.exit(1);
});
