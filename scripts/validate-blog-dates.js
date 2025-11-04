#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { execSync } from 'child_process';
import matter from 'gray-matter';

/**
 * Validates that content files don't have future dates
 */
function validateContentDates() {
  const currentDate = new Date().toISOString().split('T')[0]; // YYYY-MM-DD format
  const isPullRequest = process.env.GITHUB_EVENT_NAME === 'pull_request';
  let changedFiles = [];
  let hasErrors = false;

  console.log(`Current date: ${currentDate}`);
  console.log('');

  // Get changed files
  if (isPullRequest) {
    try {
      const output = execSync('git diff --name-only origin/main...HEAD', { encoding: 'utf8' });
      changedFiles = output
        .trim()
        .split('\n')
        .filter(file => file.startsWith('lume/src/') && (file.endsWith('.md') || file.endsWith('.mdx')));
    } catch (error) {
      console.error('Error getting changed files:', error.message);
      process.exit(1);
    }
  } else {
    // For pushes to main, get all content files
    try {
      const output = execSync('find lume/src -name "*.md" -o -name "*.mdx"', { encoding: 'utf8' });
      changedFiles = output.trim().split('\n').filter(file => file);
    } catch (error) {
      console.error('Error finding content files:', error.message);
      process.exit(1);
    }
  }

  if (changedFiles.length === 0) {
    console.log('No content files changed');
    return;
  }

  console.log(`Checking ${changedFiles.length} content file(s)...`);
  console.log('');

  for (const file of changedFiles) {
    if (!fs.existsSync(file)) continue;

    try {
      const fileContent = fs.readFileSync(file, 'utf8');
      const { data: frontmatter } = matter(fileContent);

      console.log(`Checking: ${file}`);

      if (!frontmatter.date) {
        console.log('  ⚠️  WARNING: No date field found in frontmatter');
        console.log('');
        continue;
      }

      const postDate = frontmatter.date;

      // Handle different date formats
      let dateToCheck;
      if (typeof postDate === 'string') {
        dateToCheck = postDate.split('T')[0]; // Remove time component if present
      } else if (postDate instanceof Date) {
        dateToCheck = postDate.toISOString().split('T')[0];
      } else {
        console.log(`  ⚠️  WARNING: Invalid date format: ${postDate}`);
        console.log('');
        continue;
      }

      console.log(`  Date: ${dateToCheck}`);

      if (dateToCheck > currentDate) {
        console.log('  ❌ ERROR: Future date detected!');
        hasErrors = true;
      } else {
        console.log('  ✅ OK: Date is not in the future');
      }
      console.log('');

    } catch (error) {
      console.log(`  ❌ ERROR: Failed to parse ${file}: ${error.message}`);
      hasErrors = true;
    }
  }

  if (hasErrors) {
    console.log('');
    console.log('❌ Validation failed: Found content with future dates!');
    console.log(`All content dates must be on or before the current date (${currentDate})`);
    console.log('');
    console.log('To fix this:');
    console.log('1. Update the "date" field in the frontmatter to today\'s date or earlier');
    console.log('2. Use format: date: YYYY-MM-DD');
    process.exit(1);
  } else {
    console.log('✅ All content dates are valid!');
  }
}

// Run the validation
validateContentDates();