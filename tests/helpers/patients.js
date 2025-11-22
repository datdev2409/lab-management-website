// Helper functions for patient operations in Playwright tests
const { expect } = require('@playwright/test');

/**
 * Navigate to patients page
 * @param {import('@playwright/test').Page} page
 */
async function goToPatients(page) {
  await page.goto('/danh-muc-benh-nhan');
  await page.waitForTimeout(1000);
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
  await page.waitForSelector('.alert-success');
}

/**
 * Search for a patient
 * @param {import('@playwright/test').Page} page
 * @param {string} searchTerm
 */
async function searchPatient(page, searchTerm) {
  await page.getByRole('textbox', { name: 'Tên bệnh nhân hoặc số điện thoại' }).fill(searchTerm);
  await page.waitForTimeout(500);
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

  await page.waitForTimeout(1000);
  
  if (newData.name) await row.getByTestId('patient-name-input').getByRole('textbox').fill(newData.name);
  if (newData.yob) await row.getByTestId('patient-yob-input').getByRole('textbox').fill(newData.yob);
  if (newData.gender) await row.getByTestId('patient-gender-input').getByRole('textbox').fill(newData.gender);
  if (newData.address) await row.getByTestId('patient-address-input').getByRole('textbox').fill(newData.address);
  if (newData.phone) await row.getByTestId('patient-phone-input').getByRole('textbox').fill(newData.phone);

  await page.click('text=Lưu');
  
  // Handle confirmation dialog if any
  page.on('dialog', dialog => dialog.accept());
  await page.waitForTimeout(500);
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
  await page.waitForTimeout(1000);
}

module.exports = {
  goToPatients,
  createPatient,
  searchPatient,
  editPatient,
  deletePatient,
};
