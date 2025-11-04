// Doctor Management E2E tests
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToDoctors, createDoctor, searchDoctor, editDoctor, deleteDoctor } = require('../helpers/doctors');

test.describe('Doctor Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await login(page);
  });

  test('should display doctor management page', async ({ page }) => {
    await goToDoctors(page);
    
    await expect(page.locator('h3, h4')).toContainText('bác sĩ');
    await expect(page.locator('text=Thêm bác sĩ')).toBeVisible();
  });

  test('should create a new doctor', async ({ page }) => {
    await goToDoctors(page);
    
    const doctorData = {
      name: `Dr. Test ${Date.now()}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
      email: `doctor${Date.now()}@test.com`,
      specialization: 'General Medicine',
    };
    
    await createDoctor(page, doctorData);
    
    // Verify doctor appears in the list
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    await expect(page.locator(`text=${doctorData.phone}`)).toBeVisible();
  });

  test('should search for doctors', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a test doctor first
    const doctorData = {
      name: `Dr. Search Test ${Date.now()}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createDoctor(page, doctorData);
    
    // Search for the doctor
    await searchDoctor(page, doctorData.name);
    
    // Verify search results
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    
    // Search for non-existent doctor
    await searchDoctor(page, 'NonExistentDoctor99999');
    await expect(page.locator('text=NonExistentDoctor99999')).not.toBeVisible();
  });

  test('should edit doctor information', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a doctor
    const originalData = {
      name: `Dr. Edit Test ${Date.now()}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createDoctor(page, originalData);
    
    // Edit the doctor
    const updatedData = {
      name: `${originalData.name} Updated`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await searchDoctor(page, originalData.name);
    await editDoctor(page, originalData.name, updatedData);
    
    // Verify updates
    await expect(page.locator(`text=${updatedData.name}`)).toBeVisible();
    await expect(page.locator(`text=${updatedData.phone}`)).toBeVisible();
  });

  test('should delete a doctor', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a doctor to delete
    const doctorData = {
      name: `Dr. Delete Test ${Date.now()}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createDoctor(page, doctorData);
    
    // Verify doctor exists
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    
    // Delete the doctor
    await searchDoctor(page, doctorData.name);
    await deleteDoctor(page, doctorData.name);
    
    // Verify doctor is removed
    await expect(page.locator(`text=${doctorData.name}`)).not.toBeVisible();
  });

  test('should validate required fields when creating doctor', async ({ page }) => {
    await goToDoctors(page);
    
    await page.click('text=Thêm bác sĩ');
    
    // Try to submit without filling required fields
    await page.click('button[type="submit"]');
    
    // Check for HTML5 validation
    const nameInput = page.locator('input[name="doctor_name"]');
    const isInvalid = await nameInput.evaluate(el => !el.checkValidity());
    expect(isInvalid).toBeTruthy();
  });
});
