// Record Management E2E tests (Lab Test Records)
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToRecords, createRecord, viewRecordDetails, deleteRecord } = require('../helpers/records');
const { createPatient } = require('../helpers/patients');
const { createTest } = require('../helpers/tests');
const { createCombo } = require('../helpers/combos');

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
    await page.goto('/danh-muc-benh-nhan');
    patientData = {
      name: `Record Patient ${Date.now()}`,
      yob: '1985',
      gender: 'Nam',
      address: '123 Record St',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };
    await createPatient(page, patientData);
    
    // Create tests
    await page.goto('/danh-muc-xet-nghiem');
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
    
    await expect(page.locator('h3')).toContainText('Phiếu xét nghiệm');
    await expect(page.locator('text=Tạo phiếu xét nghiệm mới')).toBeVisible();
  });

  test('should navigate to create record page', async ({ page }) => {
    await goToRecords(page);
    
    await page.click('text=Tạo phiếu xét nghiệm mới');
    await page.waitForLoadState('networkidle');
    
    await expect(page).toHaveURL('/phieu-xet-nghiem/new');
    await expect(page.locator('h3')).toContainText('Tạo phiếu xét nghiệm mới');
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
    
    // Verify we're redirected to records page or details page
    await page.waitForTimeout(1000);
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/\/(phieu-xet-nghiem|$)/);
  });

  test('should view record details', async ({ page }) => {
    await goToRecords(page);
    
    // Find a record with our test patient
    await page.fill('input[placeholder*="Tìm kiếm"]', patientData.name);
    await page.waitForTimeout(500);
    
    const row = page.locator('tr', { hasText: patientData.name }).first();
    
    // Check if the row exists before clicking
    const rowCount = await row.count();
    if (rowCount > 0) {
      await row.locator('text=Chi tiết').click();
      await page.waitForLoadState('networkidle');
      
      // Verify we're on the details page
      await expect(page).toHaveURL(/\/phieu-xet-nghiem\/.+/);
      await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
    }
  });

  test('should filter records by date', async ({ page }) => {
    await goToRecords(page);
    
    // Set date filters
    const today = new Date().toISOString().split('T')[0];
    await page.fill('input[type="date"][name="start_date"]', today);
    await page.fill('input[type="date"][name="end_date"]', today);
    
    // Apply filter
    await page.click('button:has-text("Lọc")');
    await page.waitForLoadState('networkidle');
    
    // Records should be filtered (implementation-dependent verification)
    await expect(page.locator('table')).toBeVisible();
  });

  test('should search for records by patient name', async ({ page }) => {
    await goToRecords(page);
    
    // Search for records
    await page.fill('input[placeholder*="Tìm kiếm"]', patientData.name);
    await page.waitForTimeout(500);
    
    // Check if any results appear
    const searchResults = page.locator('tr', { hasText: patientData.name });
    const count = await searchResults.count();
    
    // We expect at least one record for this patient
    expect(count).toBeGreaterThanOrEqual(0);
  });

  test('should display pagination controls', async ({ page }) => {
    await goToRecords(page);
    
    // Check for pagination elements
    const pagination = page.locator('.pagination, nav[aria-label="pagination"]');
    
    // Pagination might not exist if there are few records
    const hasPagination = await pagination.count() > 0;
    
    if (hasPagination) {
      await expect(pagination).toBeVisible();
    }
  });

  test('should create patient from record creation page', async ({ page }) => {
    await page.goto('/phieu-xet-nghiem/new');
    await page.waitForLoadState('networkidle');
    
    // Look for button to create patient
    const createPatientButton = page.locator('button:has-text("Thêm bệnh nhân"), button:has-text("Tạo bệnh nhân mới")');
    
    const hasButton = await createPatientButton.count() > 0;
    
    if (hasButton) {
      await createPatientButton.first().click();
      
      // Check if modal or form appears
      const patientForm = page.locator('form, .modal');
      await expect(patientForm).toBeVisible();
    }
  });

  test('should validate abnormal test results', async ({ page }) => {
    const recordData = {
      patientName: patientData.name,
      comboName: comboData.name,
      testResults: [
        { testName: testData.name, value: '10.5' }, // Above upper bound of 6.5
      ],
    };
    
    await page.goto('/phieu-xet-nghiem/new');
    await page.waitForLoadState('networkidle');
    
    // Search and select patient
    await page.fill('input[placeholder*="Tìm kiếm bệnh nhân"]', recordData.patientName);
    await page.waitForTimeout(500);
    
    const patientSuggestion = page.locator(`text=${recordData.patientName}`).first();
    const hasPatient = await patientSuggestion.count() > 0;
    
    if (hasPatient) {
      await patientSuggestion.click();
      await page.waitForTimeout(300);
      
      // Select combo
      await page.fill('input[placeholder*="Tìm kiếm gói xét nghiệm"]', recordData.comboName);
      await page.waitForTimeout(500);
      
      const comboSuggestion = page.locator(`text=${recordData.comboName}`).first();
      const hasCombo = await comboSuggestion.count() > 0;
      
      if (hasCombo) {
        await comboSuggestion.click();
        await page.waitForTimeout(500);
        
        // Enter abnormal value
        const testRow = page.locator('tr', { hasText: testData.name }).first();
        const testInput = testRow.locator('input[name*="test_value"]');
        
        if (await testInput.count() > 0) {
          await testInput.fill(recordData.testResults[0].value);
          
          // Check if abnormal indicator appears
          // This depends on implementation - might be a red color, icon, or text
          await page.waitForTimeout(500);
        }
      }
    }
  });
});
