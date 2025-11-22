// Helper functions for doctor operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to doctors page
 * @param {import('@playwright/test').Page} page
 */
async function goToDoctors(page) {
  await page.goto('/danh-muc-bac-si');
  await page.waitForTimeout(1000);
}

/**
 * Create a new doctor
 * @param {import('@playwright/test').Page} page
 * @param {Object} doctorData
 * @param {string} doctorData.name - Doctor name
 * @param {string} doctorData.phone - Phone number
 * @param {string} [doctorData.address] - Address
 */
async function createDoctor(page, doctorData) {
  // Open modal
  await page.getByRole('button', { name: /Thêm bác sĩ/ }).first().click();
  await page.waitForTimeout(300);
  
  // Fill form fields using specific IDs to avoid strict mode violations
  await page.locator('#doctor_name-input').fill(doctorData.name);
  await page.locator('input[name="doctor_phone"]').first().fill(doctorData.phone);
  
  if (doctorData.address) {
    await page.locator('input[name="doctor_address"]').first().fill(doctorData.address);
  }
  
  // Submit form - get button inside modal
  await page.locator('#doctor_create_form').getByRole('button', { name: /Thêm bác sĩ/ }).click();
  
  // Wait for modal to close and data to load
  await page.waitForTimeout(500);
  await page.waitForTimeout(1000);
}

/**
 * Search for a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchDoctor(page, searchTerm) {
  // Clear previous search
  const searchInput = page.getByPlaceholder(/Tìm kiếm theo tên, số điện thoại/);
  await searchInput.clear();
  
  // Fill search term - Alpine.js uses debounced input
  await searchInput.fill(searchTerm);
  
  // Wait for Alpine.js debounce (300ms) + API response
  await page.waitForTimeout(500);
  await page.waitForTimeout(1000);
}

/**
 * Edit a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} doctorName - Name of doctor to edit
 * @param {Object} newData - New doctor data
 */
async function editDoctor(page, doctorName, newData) {
  // Find the doctor row and click edit button
  const row = page.locator('tr').filter({ hasText: doctorName }).first();
  
  // Click edit button (pencil icon) to enter edit mode
  await row.locator('button', { hasText: /Sửa/ }).click();
  await page.waitForTimeout(300);
  
  // Update fields within the row
  if (newData.name) {
    const nameInput = row.locator('input[name="doctor_name"]');
    await nameInput.clear();
    await nameInput.fill(newData.name);
  }
  
  if (newData.phone) {
    const phoneInput = row.locator('input[name="doctor_phone"]');
    await phoneInput.clear();
    await phoneInput.fill(newData.phone);
  }
  
  if (newData.address) {
    const addressInput = row.locator('input[name="doctor_address"]');
    await addressInput.clear();
    await addressInput.fill(newData.address);
  }
  
  // Click save button - this triggers the API call and page redirect
  await row.locator('button:has-text("Lưu")').click();
  
  // Wait for the page to start navigating (redirect happens)
  await page.waitForNavigation({ waitUntil: 'networkidle', timeout: 10000 }).catch(() => {
    // Sometimes redirect doesn't trigger full navigation, so catch and continue
  });
  
  // Additional wait to ensure page is fully loaded
  await page.waitForTimeout(1000);
  await page.waitForTimeout(1000);
}

/**
 * Delete a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} doctorName - Name of doctor to delete
 */
async function deleteDoctor(page, doctorName) {
  const row = page.locator('tr').filter({ hasText: doctorName }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  // Click delete button (trash icon with "Xóa" text)
  await row.locator('button', { hasText: /Xóa/ }).click();
  await page.waitForTimeout(500);
  await page.waitForTimeout(1000);
}

module.exports = {
  goToDoctors,
  createDoctor,
  searchDoctor,
  editDoctor,
  deleteDoctor,
};
