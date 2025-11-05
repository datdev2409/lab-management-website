// Patient Management E2E tests
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToPatients, createPatient, searchPatient, editPatient, deletePatient } = require('../helpers/patients');

test.describe('Patient Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await login(page);
  });

  test('should display patient management page', async ({ page }) => {
    await goToPatients(page);
    
    await expect(page.getByRole('heading', { name: 'Danh mục bệnh nhân' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Thêm bệnh nhân' })).toBeVisible();
    await expect(page.getByPlaceholder('Tên bệnh nhân hoặc số điện thoại')).toBeVisible();
  });

  test('should create a new patient', async ({ page }) => {
    await goToPatients(page);
    
    const patientData = {
      name: `Test Patient ${Date.now()}`,
      yob: '1990',
      gender: 'Nam',
      address: '123 Test Street, Test City',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createPatient(page, patientData);

    await searchPatient(page, patientData.name);
    
    // Verify patient appears in the list
    await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
    await expect(page.locator(`text=${patientData.phone}`)).toBeVisible();
  });

  test('should search for patients', async ({ page }) => {
    await goToPatients(page);
    
    // Create a test patient first
    const patientData = {
      name: `Search Test ${Date.now()}`,
      yob: '1985',
      gender: 'Nữ',
      address: '456 Search Ave',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createPatient(page, patientData);
    
    // Search for the patient
    await searchPatient(page, patientData.name);
    
    // Verify search results
    await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
    
    // Search for non-existent patient
    await searchPatient(page, 'NonExistentPatient12345');
    await expect(page.locator('text=NonExistentPatient12345')).not.toBeVisible();
  });

  test('should edit patient information', async ({ page }) => {
    await goToPatients(page);
    
    // Create a patient
    const originalData = {
      name: `Edit Test ${Date.now()}`,
      yob: '1988',
      gender: 'Nam',
      address: '789 Original St',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createPatient(page, originalData);
    
    // Edit the patient
    const updatedData = {
      name: `${originalData.name} Updated`,
      address: '789 Updated Avenue',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await searchPatient(page, originalData.name);
    await editPatient(page, originalData.name, updatedData);
    await searchPatient(page, updatedData.name);
    
    // Verify updates
    await expect(page.locator(`text=${updatedData.name}`)).toBeVisible();
    await expect(page.locator(`text=${updatedData.address}`)).toBeVisible();
    await expect(page.locator(`text=${updatedData.phone}`)).toBeVisible();
  });

  test('should delete a patient', async ({ page }) => {
    await goToPatients(page);
    
    // Create a patient to delete
    const patientData = {
      name: `Delete Test ${Date.now()}`,
      yob: '1992',
      gender: 'Nữ',
      address: '321 Delete Road',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createPatient(page, patientData);

    // Verify patient exists
    await searchPatient(page, patientData.name);
    await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
    
    // Delete the patient
    await searchPatient(page, patientData.name);
    await deletePatient(page, patientData.name);
    
    // Verify patient is removed
    await searchPatient(page, patientData.name);
    await expect(page.locator(`text=${patientData.name}`)).not.toBeVisible();
  });

  test('should validate required fields when creating patient', async ({ page }) => {
    await goToPatients(page);
    
    await page.click('text=Thêm bệnh nhân');
    
    // Try to submit without filling required fields
    await page.click('button[type="submit"]');
    
    // Check for HTML5 validation
    const nameInput = page.getByRole('textbox', { name: 'Bệnh nhân', exact: true });
    const isInvalid = await nameInput.evaluate(el => !el.checkValidity());
    expect(isInvalid).toBeTruthy();
  });
});
