// Helper functions for authentication in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Login to the application
 * @param {import('@playwright/test').Page} page
 * @param {string} username
 * @param {string} password
 */
async function login(page, username = 'admin', password = 'admin123') {
  await page.goto('/login');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="password"]', password);
  await page.click('button[type="submit"]');
  
  // Wait for navigation to complete
  await page.waitForURL('/');
}

/**
 * Register a new user
 * @param {import('@playwright/test').Page} page
 * @param {string} username
 * @param {string} password
 */
async function register(page, username, password) {
  await page.goto('/register');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="password"]', password);
  await page.fill('input[name="confirm_password"]', password);
  await page.click('button[type="submit"]');
  
  // Wait for successful registration (redirects to login or home)
  await page.waitForLoadState('networkidle');
}

/**
 * Logout from the application
 * @param {import('@playwright/test').Page} page
 */
async function logout(page) {
  // Find and click logout button/link
  await page.click('text=Đăng xuất');
  await page.waitForURL('/login');
}

module.exports = {
  login,
  register,
  logout,
};
