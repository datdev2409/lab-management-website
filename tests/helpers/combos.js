// Helper functions for combo (test package) operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to combos page
 * @param {import('@playwright/test').Page} page
 */
async function goToCombos(page) {
  await page.goto('/danh-muc-goi-xet-nghiem');
  await page.waitForLoadState('networkidle');
}

/**
 * Create a new combo
 * @param {import('@playwright/test').Page} page
 * @param {Object} comboData
 * @param {string} comboData.name - Combo name
 * @param {string[]} comboData.tests - Array of test names to include
 */
async function createCombo(page, comboData) {
  await page.goto('/danh-muc-goi-xet-nghiem/new');
  await page.waitForLoadState('networkidle');
  
  await page.fill('input[name="combo_name"]', comboData.name);
  
  // Select tests for the combo
  for (const testName of comboData.tests) {
    // Type in search to find test
    await page.fill('input[placeholder*="Tìm kiếm xét nghiệm"]', testName);
    await page.waitForTimeout(500); // Wait for autocomplete
    
    // Click on the test in the dropdown
    await page.click(`text=${testName}`);
    await page.waitForTimeout(300);
  }
  
  await page.click('button[type="submit"]:has-text("Tạo gói xét nghiệm")');
  await page.waitForLoadState('networkidle');
}

/**
 * Search for a combo
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchCombo(page, searchTerm) {
  await page.fill('#combo-search', searchTerm);
  await page.click('#combo-search-form button[type="submit"]');
  await page.waitForLoadState('networkidle');
}

/**
 * Delete a combo
 * @param {import('@playwright/test').Page} page
 * @param {string} comboName - Name of combo to delete
 */
async function deleteCombo(page, comboName) {
  const row = page.locator('tr', { hasText: comboName }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  await row.locator('text=Xoá').click();
  await page.waitForLoadState('networkidle');
}

module.exports = {
  goToCombos,
  createCombo,
  searchCombo,
  deleteCombo,
};
