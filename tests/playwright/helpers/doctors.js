// Helper functions for doctor operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to doctors page
 * @param {import('@playwright/test').Page} page
 */
async function goToDoctors(page) {
  await page.goto('/danh-muc-bac-si');
  await page.waitForLoadState('networkidle');
}

/**
 * Create a new doctor
 * @param {import('@playwright/test').Page} page
 * @param {Object} doctorData
 * @param {string} doctorData.name - Doctor name
 * @param {string} doctorData.phone - Phone number
 * @param {string} [doctorData.email] - Email address
 * @param {string} [doctorData.specialization] - Specialization
 */
async function createDoctor(page, doctorData) {
  await page.click('text=Thêm bác sĩ');
  
  await page.fill('input[name="doctor_name"]', doctorData.name);
  await page.fill('input[name="doctor_phone"]', doctorData.phone);
  
  if (doctorData.email) {
    await page.fill('input[name="doctor_email"]', doctorData.email);
  }
  
  if (doctorData.specialization) {
    await page.fill('input[name="doctor_specialization"]', doctorData.specialization);
  }
  
  await page.click('button[type="submit"]');
  await page.waitForLoadState('networkidle');
}

/**
 * Search for a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchDoctor(page, searchTerm) {
  await page.fill('#doctor-search', searchTerm);
  await page.click('#doctor-search-form button[type="submit"]');
  await page.waitForLoadState('networkidle');
}

/**
 * Edit a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} doctorName - Name of doctor to edit
 * @param {Object} newData - New doctor data
 */
async function editDoctor(page, doctorName, newData) {
  const row = page.locator('tr', { hasText: doctorName }).first();
  await row.locator('text=Sửa').click();
  
  if (newData.name) await page.fill('input[name="doctor_name"]', newData.name);
  if (newData.phone) await page.fill('input[name="doctor_phone"]', newData.phone);
  if (newData.email) await page.fill('input[name="doctor_email"]', newData.email);
  if (newData.specialization) await page.fill('input[name="doctor_specialization"]', newData.specialization);
  
  await page.click('text=Lưu');
  
  // Handle confirmation dialog if any
  page.on('dialog', dialog => dialog.accept());
  await page.waitForLoadState('networkidle');
}

/**
 * Delete a doctor
 * @param {import('@playwright/test').Page} page
 * @param {string} doctorName - Name of doctor to delete
 */
async function deleteDoctor(page, doctorName) {
  const row = page.locator('tr', { hasText: doctorName }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  await row.locator('text=Xoá').click();
  await page.waitForLoadState('networkidle');
}

module.exports = {
  goToDoctors,
  createDoctor,
  searchDoctor,
  editDoctor,
  deleteDoctor,
};
