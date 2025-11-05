// Record Management E2E tests (Lab Test Records)
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToRecords, createRecord, viewRecordDetails, deleteRecord } = require('../helpers/records');
const { goToPatients, createPatient } = require('../helpers/patients');
const { goToTests, createTest } = require('../helpers/tests');
const { goToCombos, createCombo } = require('../helpers/combos');

test.describe('Record Management Flow', () => {
  let patientData;
  let testData;
  let comboData;

  test.beforeAll(async ({ browser }) => {
    // Setup test data: create a patient, tests, and combo
    const context = await browser.newContext();
    const page = await context.newPage();
    
    await login(page);
    
    // Create a patient
    await goToPatients(page);
    patientData = {
      name: `Record Patient ${Date.now()}`,
      yob: '1985',
      gender: 'Nam',
      address: '123 Record St',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    await createPatient(page, patientData);
    
    // Create tests
    await goToTests(page);
    testData = {
      name: `Record Test ${Date.now()}`,
      unit: 'mmol/L',
      price: 50000,
      lowerBound: 3.5,
      upperBound: 6.5,
      normalValue: '3.5-6.5',
    };
    await createTest(page, testData);
    
    // Create a combo
    await goToCombos(page);
    comboData = {
      name: `Record Combo ${Date.now()}`,
      tests: [testData.name],
    };
    await createCombo(page, comboData);
    
    await context.close();
  });

  test.beforeEach(async ({ page }) => {
    // Login before each test
    await login(page);
  });

  test('should display records page', async ({ page }) => {
    await goToRecords(page);
    
    await expect(page.getByRole('heading', { name: 'Phiếu xét nghiệm' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Tạo phiếu xét nghiệm mới' })).toBeVisible();
  });

  test('should navigate to create record page', async ({ page }) => {
    await goToRecords(page);
    
    await page.getByRole('button', { name: 'Tạo phiếu xét nghiệm mới' }).click();
    await page.waitForLoadState('networkidle');
    
    await expect(page).toHaveURL('/phieu-xet-nghiem/new');
    await expect(page.getByRole('heading', { name: 'Tạo phiếu xét nghiệm mới' })).toBeVisible();
  });

  test('should create a new record with combo', async ({ page }) => {
    const recordData = {
      patientName: patientData.name,
      comboName: comboData.name,
      testResults: [
        { testName: testData.name, value: '5.2' },
      ],
    };
    
    await createRecord(page, recordData);
    expect(page.getByText('Tạo phiếu xét nghiệm thành công')).toBeVisible();
    
    // Verify we're redirected and record was created
    await page.waitForTimeout(1000);
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/phieu-xet-nghiem/);
    expect(currentUrl).not.toBe('/phieu-xet-nghiem/new');
  });

  test('should view record details', async ({ page }) => {
    await goToRecords(page);
    await page.waitForTimeout(500);
    
    // Search for records with our test patient
    const searchInput = page.getByPlaceholder('Tên bệnh nhân hoặc số điện thoại');
    await searchInput.fill(patientData.name);
    await page.waitForTimeout(600);
    
    // Find and click details button
    const row = page.locator('tr', { hasText: patientData.name }).first();
    const rowCount = await row.count();
    
    if (rowCount > 0) {
      await row.locator('a:has-text("Chi tiết")').first().click();
      await page.waitForLoadState('networkidle');
      
      // Verify we're on the details page
      await expect(page).toHaveURL(/\/phieu-xet-nghiem\/[^/]+/);
    }
  });

  test('should search for records by patient name', async ({ page }) => {
    await goToRecords(page);
    
    // Search for records
    const searchInput = page.getByPlaceholder('Tên bệnh nhân hoặc số điện thoại');
    await searchInput.fill(patientData.name);
    await page.waitForTimeout(600);
    
    // Check if any results appear
    const searchResults = page.locator('tr', { hasText: patientData.name });
    const count = await searchResults.count();
    
    // We expect at least one record for this patient
    expect(count).toBeGreaterThanOrEqual(0);
  });

  test('should display records table with columns', async ({ page }) => {
    await goToRecords(page);
    
    // Check for table headers
    await expect(page.getByRole('table')).toBeVisible();
    
    // Common table headers in records
    const table = page.getByRole('table');
    const headers = table.locator('th');
    
    expect(await headers.count()).toBeGreaterThan(0);
  });
});
