#!/usr/bin/env node

/**
 * Extract title and publication date from HTML HEAD section
 * Usage: node scripts/extract-meta.js <url>
 * Example: node scripts/extract-meta.js https://www.tigrisdata.com/blog/storage-sdk/
 */

import https from 'https';
import http from 'http';
import { JSDOM } from 'jsdom';

// Function to fetch HTML content
function fetchHTML(url) {
  return new Promise((resolve, reject) => {
    const client = url.startsWith('https:') ? https : http;

    const request = client.get(url, (response) => {
      let html = '';

      // Handle redirects
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        return fetchHTML(response.headers.location).then(resolve).catch(reject);
      }

      response.setEncoding('utf8');
      response.on('data', (chunk) => {
        html += chunk;
      });

      response.on('end', () => {
        if (response.statusCode >= 200 && response.statusCode < 300) {
          resolve(html);
        } else {
          reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
        }
      });
    });

    request.on('error', (err) => {
      reject(err);
    });

    request.setTimeout(10000, () => {
      request.destroy();
      reject(new Error('Request timeout'));
    });
  });
}

// Function to extract publication date from various meta tags
function extractPublicationDate(document) {
  const dateSelectors = [
    'meta[property="article:published_time"]',
    'meta[property="article:published"]',
    'meta[name="article:published_time"]',
    'meta[name="publication_date"]',
    'meta[name="date"]',
    'meta[property="datePublished"]',
    'meta[name="DC.date"]',
    'meta[name="DC.date.created"]',
    'time[datetime]'
  ];

  for (const selector of dateSelectors) {
    const element = document.querySelector(selector);
    if (element) {
      const date = element.getAttribute('content') || element.getAttribute('datetime');
      if (date) {
        return parseDate(date);
      }
    }
  }

  return null;
}

// Function to parse and normalize date
function parseDate(dateString) {
  try {
    // Handle ISO 8601 dates
    if (dateString.includes('T') || dateString.includes('Z')) {
      const date = new Date(dateString);
      if (!isNaN(date.getTime())) {
        return date.toISOString().split('T')[0]; // Return YYYY-MM-DD
      }
    }

    // Handle various date formats
    const date = new Date(dateString);
    if (!isNaN(date.getTime())) {
      return date.toISOString().split('T')[0];
    }

    // Try to extract date from string patterns
    const dateMatch = dateString.match(/(\d{4})-(\d{2})-(\d{2})/);
    if (dateMatch) {
      return dateMatch[0];
    }

    return dateString; // Return original if parsing fails
  } catch (error) {
    console.warn(`Failed to parse date: ${dateString}`, error.message);
    return dateString;
  }
}

// Main function
async function main() {
  const url = process.argv[2];

  if (!url) {
    console.error('Usage: node extract-meta.js <url>');
    console.error('Example: node extract-meta.js https://www.tigrisdata.com/blog/storage-sdk/');
    process.exit(1);
  }

  try {
    console.log(`Fetching: ${url}`);
    const html = await fetchHTML(url);

    const dom = new JSDOM(html);
    const document = dom.window.document;

    // Extract title
    const title = document.querySelector('title')?.textContent?.trim() || '';

    // Extract Open Graph title as fallback
    const ogTitle = document.querySelector('meta[property="og:title"]')?.getAttribute('content')?.trim() || '';

    // Extract publication date
    const publicationDate = extractPublicationDate(document);

    // Output results as JSON
    const result = {
      url,
      title: title || ogTitle,
      publicationDate
    };

    console.log(JSON.stringify(result, null, 2));

  } catch (error) {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  }
}

main();