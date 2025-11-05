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
    
    // Verify page heading
    await expect(page.getByRole('heading', { name: /Danh mục bác sĩ/ })).toBeVisible();
    
    // Verify "Add doctor" button
    await expect(page.getByRole('button', { name: /Thêm bác sĩ/ })).toBeVisible();
    
    // Verify search input
    await expect(page.getByPlaceholder(/Tìm kiếm theo tên, số điện thoại/)).toBeVisible();
  });

  test('should create a new doctor', async ({ page }) => {
    await goToDoctors(page);
    
    const timestamp = Date.now();
    const doctorData = {
      name: `Dr. Test ${timestamp}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
      address: `Test Address ${timestamp}`,
    };
    
    await createDoctor(page, doctorData);
    
    // Verify doctor appears in the list
    await searchDoctor(page, doctorData.name);
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    await expect(page.locator(`text=${doctorData.phone}`)).toBeVisible();
  });

  test('should search for doctors', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a test doctor first
    const timestamp = Date.now();
    const doctorData = {
      name: `Dr. Search Test ${timestamp}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createDoctor(page, doctorData);
    
    // Clear previous search and search for the doctor
    await searchDoctor(page, doctorData.name);
    
    // Verify search results
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    
    // Search for non-existent doctor
    await searchDoctor(page, 'NonExistentDoctor99999');
    
    // Verify "not found" message appears instead of results
    await expect(page.locator('text=Không tìm thấy bác sĩ nào')).toBeVisible();
  });

  test('should edit doctor information', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a doctor
    const timestamp = Date.now();
    const originalData = {
      name: `Dr. Edit Test ${timestamp}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
      address: 'Original Address',
    };
    
    await createDoctor(page, originalData);
    
    // Verify the original doctor is in the list
    await searchDoctor(page, originalData.name);
    await expect(page.locator(`text=${originalData.name}`)).toBeVisible();
    
    // Prepare updated data
    const updatedData = {
      name: `Dr. Edit Test ${timestamp} Updated`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
      address: 'Updated Address',
    };
    
    // Edit the doctor (using original name to find the row)
    await editDoctor(page, originalData.name, updatedData);
    await searchDoctor(page, updatedData.name);
    
    // After redirect, verify updates appear in the list
    await expect(page.locator(`text=${updatedData.name}`)).toBeVisible({ timeout: 8000 });
    await expect(page.locator(`text=${updatedData.phone}`)).toBeVisible();
  });

  test('should delete a doctor', async ({ page }) => {
    await goToDoctors(page);
    
    // Create a doctor to delete
    const timestamp = Date.now();
    const doctorData = {
      name: `Dr. Delete Test ${timestamp}`,
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    
    await createDoctor(page, doctorData);
    
    // Search for the doctor to ensure it exists
    await searchDoctor(page, doctorData.name);
    await expect(page.locator(`text=${doctorData.name}`)).toBeVisible();
    
    // Delete the doctor
    await deleteDoctor(page, doctorData.name);

    await searchDoctor(page, doctorData.name);
    
    // Verify doctor is removed
    await expect(page.locator('text=Không tìm thấy bác sĩ nào')).toBeVisible();
    await expect(page.locator(`text=${doctorData.name}`)).not.toBeVisible();
  });

  test('should validate required fields when creating doctor', async ({ page }) => {
    await goToDoctors(page);
    
    // Open the create doctor modal
    await page.getByRole('button', { name: /Thêm bác sĩ/ }).click();
    await page.waitForTimeout(300);
    
    // Try to submit without filling required fields
    const submitButton = page.getByRole('button', { name: /Thêm bác sĩ/ }).last();
    
    // Click submit - Alpine validation should show error
    await page.locator('#doctor_create_form').evaluate((form) => {
      form.dispatchEvent(new Event('submit'));
    });
    
    // Wait for validation error to appear
    await page.waitForTimeout(500);
    
    // Check if error alert appears
    const errorAlert = page.locator('.alert-danger', { hasText: /không được để trống/ });
    const hasError = await errorAlert.count() > 0 || (await submitButton.evaluate((btn) => btn.disabled));
    
    expect(hasError).toBeTruthy();
  });

  test('should display multiple doctors in the table', async ({ page }) => {
    await goToDoctors(page);
    
    const timestamp = Date.now();
    const doctors = [];
    
    // Create 3 doctors
    for (let i = 1; i <= 3; i++) {
      const doctorData = {
        name: `Dr. Multi Test ${timestamp} - ${i}`,
        phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
      };
      doctors.push(doctorData);
      await createDoctor(page, doctorData);
    }
    
    // Clear search to see all doctors
    await searchDoctor(page, `Dr. Multi Test ${timestamp}`);
    await page.waitForTimeout(500);
    
    // Verify all doctors appear in the list
    for (const doctor of doctors) {
      await expect(page.locator(`text=${doctor.name}`)).toBeVisible();
    }
  });

  test('should navigate back from doctor list', async ({ page }) => {
    await goToDoctors(page);
    
    // Verify we're on the doctors page
    await expect(page.getByRole('heading', { name: /Danh mục bác sĩ/ })).toBeVisible();
    
    // Check URL
    await expect(page).toHaveURL('/danh-muc-bac-si');
  });
});
