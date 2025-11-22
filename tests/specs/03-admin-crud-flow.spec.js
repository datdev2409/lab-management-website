// Flow 2: Admin CRUD - Management of Patients, Tests, and Combos
// Consolidated from individual CRUD test files
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToPatients, createPatient, searchPatient, editPatient, deletePatient } = require('../helpers/patients');
const { goToTests, createTest, searchTest, deleteTest } = require('../helpers/tests');
const { goToCombos, createCombo, searchCombo, deleteCombo } = require('../helpers/combos');

test.describe('Flow 2: Admin CRUD - Entity Management', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test.describe('Patient Management', () => {
    test('should display patient management page', async ({ page }) => {
      await goToPatients(page);
      
      await expect(page.getByRole('heading', { name: 'Danh mục bệnh nhân' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Thêm bệnh nhân' })).toBeVisible();
      await expect(page.getByPlaceholder('Tên bệnh nhân hoặc số điện thoại')).toBeVisible();
    });

    test('should create a new patient', async ({ page }) => {
      await goToPatients(page);
      
      const patientData = {
        name: `CRUD Patient ${Date.now()}`,
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
      await searchPatient(page, patientData.name);
      await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
      
      // Delete the patient
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

  test.describe('Test Management', () => {
    test('should display test management page', async ({ page }) => {
      await goToTests(page);

      await expect(page.getByRole('heading', { name: 'Danh mục xét nghiệm' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Thêm xét nghiệm' })).toBeVisible();
      await expect(page.getByPlaceholder('Tên xét nghiệm')).toBeVisible();
    });

    test('should create a new test', async ({ page }) => {
      await goToTests(page);
      
      const testData = {
        name: `CRUD Test ${Date.now()}`,
        unit: 'mmol/L',
        price: 50000,
        lowerBound: 3.9,
        upperBound: 6.1,
        normalValue: '3.9-6.1',
      };
      
      await createTest(page, testData);
      await searchTest(page, testData.name);
      
      // Verify test appears in the list
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
      await expect(page.locator(`text=${testData.unit}`)).toBeVisible();
    });

    test('should create multiple tests with different parameters', async ({ page }) => {
      await goToTests(page);
      
      const timestamp = Date.now();
      const tests = [
        {
          name: `Multi Test A ${timestamp}`,
          unit: 'g/dL',
          price: 30000,
          lowerBound: 12.0,
          upperBound: 16.0,
          normalValue: '12.0-16.0',
        },
        {
          name: `Multi Test B ${timestamp}`,
          unit: '10^9/L',
          price: 40000,
          lowerBound: 4.0,
          upperBound: 11.0,
          normalValue: '4.0-11.0',
        },
      ];
      
      for (const testData of tests) {
        await createTest(page, testData);
        await searchTest(page, testData.name);
        await expect(page.locator(`text=${testData.name}`)).toBeVisible();
      }
    });

    test('should search for tests', async ({ page }) => {
      await goToTests(page);
      
      // Create a test first
      const testData = {
        name: `Search Test ${Date.now()}`,
        unit: 'mmol/L',
        price: 60000,
        lowerBound: 3.0,
        upperBound: 5.2,
        normalValue: '3.0-5.2',
      };
      
      await createTest(page, testData);
      await searchTest(page, testData.name);
      
      // Verify search results
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
      
      // Search for non-existent test
      await searchTest(page, 'NonExistentTest99999');
      await expect(page.locator('text=NonExistentTest99999')).not.toBeVisible();
    });

    test('should delete a test', async ({ page }) => {
      await goToTests(page);
      
      // Create a test to delete
      const testData = {
        name: `Delete Test ${Date.now()}`,
        unit: 'mg/dL',
        price: 45000,
        lowerBound: 70,
        upperBound: 100,
        normalValue: '70-100',
      };
      
      await createTest(page, testData);
      await searchTest(page, testData.name);
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
      
      // Delete the test
      await deleteTest(page, testData.name);
      
      // Verify test is removed
      await searchTest(page, testData.name);
      await expect(page.locator(`text=${testData.name}`)).not.toBeVisible();
    });

    test('should display test details correctly', async ({ page }) => {
      await goToTests(page);
      
      const testData = {
        name: `Detail Test ${Date.now()}`,
        unit: 'IU/L',
        price: 75000,
        lowerBound: 10,
        upperBound: 40,
        normalValue: '10-40',
      };
      
      await createTest(page, testData);
      await searchTest(page, testData.name);
      
      // Verify all details are displayed
      const row = page.locator('tr', { hasText: testData.name }).first();
      await expect(row).toContainText(testData.unit);
      await expect(row).toContainText('75,000'); // Price formatted
      await expect(row).toContainText(testData.normalValue);
    });

    test('should validate required fields when creating test', async ({ page }) => {
      await goToTests(page);
      
      await page.click('text=Thêm xét nghiệm');
      
      // Try to submit without filling required fields
      await page.click('button[type="submit"]:has-text("Thêm xét nghiệm")');
      
      // Check for HTML5 validation
      const nameInput = page.getByLabel('Tên xét nghiệm');
      const isInvalid = await nameInput.evaluate(el => !el.checkValidity());
      expect(isInvalid).toBeTruthy();
    });
  });

  test.describe('Combo Management', () => {
    let testNames = [];

    test.beforeAll(async ({ browser }) => {
      // Create some tests to use in combos
      const context = await browser.newContext();
      const page = await context.newPage();
      
      await login(page);
      await goToTests(page);
      
      const timestamp = Date.now();
      const tests = [
        {
          name: `Combo Test A ${timestamp}`,
          unit: 'mmol/L',
          price: 50000,
          lowerBound: 3.0,
          upperBound: 6.0,
          normalValue: '3.0-6.0',
        },
        {
          name: `Combo Test B ${timestamp}`,
          unit: 'g/dL',
          price: 40000,
          lowerBound: 12.0,
          upperBound: 16.0,
          normalValue: '12.0-16.0',
        },
      ];
      
      for (const testData of tests) {
        await createTest(page, testData);
        testNames.push(testData.name);
      }
      
      await context.close();
    });

    test('should display combo management page', async ({ page }) => {
      await goToCombos(page);
      
      await expect(page.getByRole('heading', { name: 'Danh mục gói xét nghiệm' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Tạo gói xét nghiệm mới' })).toBeVisible();
    });

    test('should create a new combo with multiple tests', async ({ page }) => {
      await goToCombos(page);
      
      const comboData = {
        name: `CRUD Combo ${Date.now()}`,
        tests: testNames,
      };
      
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      
      // Verify combo appears in the list
      await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
    });

    test('should search for combos', async ({ page }) => {
      await goToCombos(page);
      
      // Create a combo first
      const comboData = {
        name: `Search Combo ${Date.now()}`,
        tests: [testNames[0]],
      };
      
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      
      // Verify search results
      await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
      
      // Search for non-existent combo
      await searchCombo(page, 'NonExistentCombo99999');
      await expect(page.locator('text=NonExistentCombo99999')).not.toBeVisible();
    });

    test('should view combo details', async ({ page }) => {
      await goToCombos(page);
      
      // Create a combo
      const comboData = {
        name: `Detail Combo ${Date.now()}`,
        tests: testNames,
      };
      
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      
      // Click to view details
      const row = page.locator('tr', { hasText: comboData.name }).first();
      await row.locator('text=Chi tiết').click();
      await page.waitForTimeout(1000);
      
      // Verify we're on the details/edit page
      await expect(page).toHaveURL(/\/danh-muc-goi-xet-nghiem\/.*\/edit/);
    });

    test('should delete a combo', async ({ page }) => {
      await goToCombos(page);
      
      // Create a combo to delete
      const comboData = {
        name: `Delete Combo ${Date.now()}`,
        tests: [testNames[0]],
      };
      
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
      
      // Delete the combo
      await deleteCombo(page, comboData.name);
      
      // Verify combo is removed
      await searchCombo(page, comboData.name);
      await expect(page.locator(`text=${comboData.name}`)).not.toBeVisible();
    });

    test('should display combo test count in list', async ({ page }) => {
      await goToCombos(page);
      
      // Create a combo with known number of tests
      const comboData = {
        name: `Count Test Combo ${Date.now()}`,
        tests: testNames,
      };
      
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      
      // Verify test count is displayed
      const row = page.locator('tr', { hasText: comboData.name }).first();
      await expect(row).toContainText(`${testNames.length} xét nghiệm`);
    });

    test('should navigate back from combo creation form', async ({ page }) => {
      await page.goto('/danh-muc-goi-xet-nghiem/new');
      await page.waitForTimeout(1000);
      
      // Click back button
      await page.getByRole('button', { name: 'Trở lại' }).click();
      await page.waitForTimeout(500);
      
      // Verify we're back on the combo list page
      await expect(page).toHaveURL('/danh-muc-goi-xet-nghiem');
    });
  });

  test('should navigate through all admin pages', async ({ page }) => {
    const pages = [
      { url: '/danh-muc-benh-nhan', name: 'Patients' },
      { url: '/danh-muc-xet-nghiem', name: 'Tests' },
      { url: '/danh-muc-goi-xet-nghiem', name: 'Combos' },
    ];

    for (const pageInfo of pages) {
      await page.goto(pageInfo.url);
      await page.waitForTimeout(1000);
      
      // Verify page loaded by checking for main content
      const table = page.getByRole('table');
      await expect(table).toBeVisible();
      
      console.log(`✓ Navigated to ${pageInfo.name} (${pageInfo.url})`);
    }

    console.log('✓✓✓ FLOW 2 (Admin CRUD) COMPLETED SUCCESSFULLY ✓✓✓');
  });
});
