// Helper functions for test management operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to tests page
 * @param {import('@playwright/test').Page} page
 */
async function goToTests(page) {
  await page.goto('/danh-muc-xet-nghiem');
  await page.waitForTimeout(1000);
}

/**
 * Create a new test
 * @param {import('@playwright/test').Page} page
 * @param {Object} testData
 * @param {string} testData.name - Test name
 * @param {string} testData.unit - Unit of measurement
 * @param {number} testData.price - Price
 * @param {number} testData.lowerBound - Lower bound
 * @param {number} testData.upperBound - Upper bound
 * @param {string} testData.normalValue - Normal value range
 */
async function createTest(page, testData) {
  await page.click('text=Thêm xét nghiệm');
  
  await page.fill('input[name="test_name"]', testData.name);
  await page.fill('input[name="test_unit"]', testData.unit);
  await page.fill('input[name="test_price"]', String(testData.price));
  await page.fill('input[name="test_lower_bound"]', String(testData.lowerBound));
  await page.fill('input[name="test_upper_bound"]', String(testData.upperBound));
  await page.fill('input[name="test_normal_value"]', testData.normalValue);
  
  await page.click('button[type="submit"]:has-text("Thêm xét nghiệm")');
  await page.waitForTimeout(500);
}

/**
 * Search for a test
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchTest(page, searchTerm) {
  await page.getByPlaceholder('Tên xét nghiệm').fill(searchTerm);
  await page.waitForTimeout(5000);
}

/**
 * Delete a test
 * @param {import('@playwright/test').Page} page
 * @param {string} testName - Name of test to delete
 */
async function deleteTest(page, testName) {
  const row = page.locator('tr', { hasText: testName }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  await row.locator('text=Xoá').click();
  await page.waitForTimeout(1000);
}

module.exports = {
  goToTests,
  createTest,
  searchTest,
  deleteTest,
};
