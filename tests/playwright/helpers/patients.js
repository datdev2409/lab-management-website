// Helper functions for patient operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to patients page
 * @param {import('@playwright/test').Page} page
 */
async function goToPatients(page) {
  await page.goto('/danh-muc-benh-nhan');
  await page.waitForLoadState('networkidle');
}

/**
 * Create a new patient
 * @param {import('@playwright/test').Page} page
 * @param {Object} patientData
 * @param {string} patientData.name - Patient name
 * @param {string} patientData.yob - Year of birth
 * @param {string} patientData.gender - Gender (Nam/Nữ)
 * @param {string} patientData.address - Address
 * @param {string} patientData.phone - Phone number
 */
async function createPatient(page, patientData) {
  await page.click('text=Thêm bệnh nhân');
  
  await page.fill('input[name="patient_name"]', patientData.name);
  await page.fill('input[name="patient_yob"]', patientData.yob);
  await page.check(`input[name="patient_gender"][value="${patientData.gender}"]`);
  await page.fill('input[name="patient_address"]', patientData.address);
  await page.fill('input[name="patient_phone"]', patientData.phone);
  
  await page.click('button[type="submit"]');
  await page.waitForLoadState('networkidle');
}

/**
 * Search for a patient
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchPatient(page, searchTerm) {
  await page.fill('#patient-search', searchTerm);
  await page.click('#patient-search-form button[type="submit"]');
  await page.waitForLoadState('networkidle');
}

/**
 * Edit a patient
 * @param {import('@playwright/test').Page} page
 * @param {string} patientName - Name of patient to edit
 * @param {Object} newData - New patient data
 */
async function editPatient(page, patientName, newData) {
  // Find the patient row and click edit
  const row = page.locator('tr', { hasText: patientName }).first();
  await row.locator('text=Sửa').click();
  
  if (newData.name) await page.fill('input[name="patient_name"]', newData.name);
  if (newData.yob) await page.fill('input[name="patient_yob"]', newData.yob);
  if (newData.gender) await page.check(`input[name="patient_gender"][value="${newData.gender}"]`);
  if (newData.address) await page.fill('input[name="patient_address"]', newData.address);
  if (newData.phone) await page.fill('input[name="patient_phone"]', newData.phone);
  
  await page.click('text=Lưu');
  
  // Handle confirmation dialog if any
  page.on('dialog', dialog => dialog.accept());
  await page.waitForLoadState('networkidle');
}

/**
 * Delete a patient
 * @param {import('@playwright/test').Page} page
 * @param {string} patientName - Name of patient to delete
 */
async function deletePatient(page, patientName) {
  const row = page.locator('tr', { hasText: patientName }).first();
  
  // Setup dialog handler before clicking delete
  page.once('dialog', dialog => {
    expect(dialog.type()).toBe('confirm');
    dialog.accept();
  });
  
  await row.locator('text=Xoá').click();
  await page.waitForLoadState('networkidle');
}

module.exports = {
  goToPatients,
  createPatient,
  searchPatient,
  editPatient,
  deletePatient,
};
