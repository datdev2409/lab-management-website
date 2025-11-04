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
  
  // Search and select patient using autocomplete
  const patientInput = page.getByRole('row', { name: 'Bệnh nhân' }).getByRole('textbox');
  await patientInput.fill(recordData.patientName);
  await page.waitForTimeout(600); // Wait for autocomplete results
  
  // Click on patient from autocomplete dropdown
  const patientOption = page.locator('.autocomplete-option', { hasText: recordData.patientName }).first();
  await patientOption.click();
  await page.waitForTimeout(300);
  
  // Select combo if provided
  if (recordData.comboName) {
    const comboInput = page.getByRole('row', { name: 'Tên gói xét nghiệm' }).getByRole('textbox');
    await comboInput.fill(recordData.comboName);
    await page.waitForTimeout(600); // Wait for autocomplete results
    
    const comboOption = page.locator('.autocomplete-option', { hasText: recordData.comboName }).first();
    await comboOption.click();
    await page.waitForTimeout(500);
  }
  
  // Enter test results
  if (recordData.testResults) {
    for (const result of recordData.testResults) {
      // Find the test result input field by looking for the test name and then the input in that row
      const testRow = page.locator('tr', { hasText: result.testName }).first();
      const testInput = testRow.locator('input[type="number"]').first();
      
      if (await testInput.count() > 0) {
        await testInput.fill(result.value);
        await page.waitForTimeout(200);
      }
    }
  }
  
  // Submit the form
  await page.getByRole('button', { name: 'Tạo (Ctrl + S)' }).click();
  await page.waitForTimeout(500);
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
