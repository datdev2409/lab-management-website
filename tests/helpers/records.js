// Helper functions for record operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to records page
 * @param {import('@playwright/test').Page} page
 */
async function goToRecords(page) {
  await page.goto('/phieu-xet-nghiem');
  await page.waitForLoadState('networkidle');
}

/**
 * Create a new record
 * @param {import('@playwright/test').Page} page
 * @param {Object} recordData
 * @param {string} recordData.patientName - Patient name to search
 * @param {string} recordData.comboName - Combo name (if using combo)
 * @param {Object[]} recordData.testResults - Array of test results
 * @param {string} recordData.testResults[].testName - Test name
 * @param {string} recordData.testResults[].value - Test value
 */
async function createRecord(page, recordData) {
  await page.goto('/phieu-xet-nghiem/new');
  await page.waitForLoadState('networkidle');
  
  // Search and select patient
  await page.fill('input[placeholder*="Tìm kiếm bệnh nhân"]', recordData.patientName);
  await page.waitForTimeout(500); // Wait for autocomplete
  await page.click(`text=${recordData.patientName}`);
  await page.waitForTimeout(300);
  
  // Select combo if provided
  if (recordData.comboName) {
    await page.fill('input[placeholder*="Tìm kiếm gói xét nghiệm"]', recordData.comboName);
    await page.waitForTimeout(500);
    await page.click(`text=${recordData.comboName}`);
    await page.waitForTimeout(300);
  }
  
  // Enter test results
  if (recordData.testResults) {
    for (const result of recordData.testResults) {
      // Find the test result input field
      const testRow = page.locator('tr', { hasText: result.testName }).first();
      await testRow.locator('input[name*="test_value"]').fill(result.value);
    }
  }
  
  await page.click('button[type="submit"]:has-text("Tạo phiếu")');
  await page.waitForLoadState('networkidle');
}

/**
 * View record details
 * @param {import('@playwright/test').Page} page
 * @param {string} recordIdentifier - Patient name or record ID
 */
async function viewRecordDetails(page, recordIdentifier) {
  await goToRecords(page);
  
  const row = page.locator('tr', { hasText: recordIdentifier }).first();
  await row.locator('text=Chi tiết').click();
  await page.waitForLoadState('networkidle');
}

/**
 * Generate report for a record
 * @param {import('@playwright/test').Page} page
 * @param {string} reportType - Type of report (e.g., 'phieu_thu', 'phieu_ket_qua')
 */
async function generateReport(page, reportType) {
  // Assuming we're on record details page
  await page.click(`button:has-text("${reportType}")`);
  
  // Wait for the report generation
  const downloadPromise = page.waitForEvent('download');
  const download = await downloadPromise;
  
  return download;
}

/**
 * Delete a record
 * @param {import('@playwright/test').Page} page
 * @param {string} recordIdentifier - Patient name or record ID
 */
async function deleteRecord(page, recordIdentifier) {
  await goToRecords(page);
  
  const row = page.locator('tr', { hasText: recordIdentifier }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  await row.locator('text=Xoá').click();
  await page.waitForLoadState('networkidle');
}

module.exports = {
  goToRecords,
  createRecord,
  viewRecordDetails,
  generateReport,
  deleteRecord,
};
