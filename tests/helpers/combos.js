// Helper functions for combo (test package) operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to combos page
 * @param {import('@playwright/test').Page} page
 */
async function goToCombos(page) {
  await page.goto('/danh-muc-goi-xet-nghiem');
  await page.waitForTimeout(1000);
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
  await page.waitForTimeout(500);

  // Fill combo name
  await page.getByRole('textbox', { name: 'Tên gói xét nghiệm' }).fill(comboData.name);
  await page.waitForTimeout(300);

  // Select tests for the combo using the autocomplete
  for (const testName of comboData.tests) {
    // Type in the test autocomplete field
    const autocompleteInput = page.locator('input#test-autocomplete');
    await autocompleteInput.fill(testName);
    await page.waitForTimeout(600); // Wait for autocomplete to load
    
    // Click on the test option in the dropdown
    const testOption = page.locator('.autocomplete-option', { hasText: testName }).first();
    await testOption.click();
    await page.waitForTimeout(300);
  }
  
  // Submit the form
  await page.getByRole('button', { name: /Tạo gói xét nghiệm|Cập nhật gói xét nghiệm/ }).click();
  await page.waitForTimeout(1000);
}

/**
 * Search for a combo
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchCombo(page, searchTerm) {
  await page.getByPlaceholder('Tên gói xét nghiệm').fill(searchTerm);
  await page.waitForTimeout(500); // Wait for debounced search
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
  
  await row.locator('text=Xóa').click();
  await page.waitForTimeout(1000);
}

module.exports = {
  goToCombos,
  createCombo,
  searchCombo,
  deleteCombo,
};
