// Complete E2E Integration Test - Full Application Flow
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToPatients, createPatient, searchPatient } = require('../helpers/patients');
const { goToTests, createTest, searchTest } = require('../helpers/tests');
const { goToCombos, createCombo, searchCombo } = require('../helpers/combos');
const { goToRecords, createRecord } = require('../helpers/records');

test.describe('Complete Application Flow Integration Test', () => {
  // Extend timeout to 120 seconds for integration tests (instead of default 30 seconds)
  // This is necessary because these tests perform multiple CRUD operations sequentially
  test.describe.configure({ timeout: 120000 });

  test('should complete full workflow: create patient, tests, combo, and record', async ({ page }) => {
    // Step 1: Login
    await login(page);
    await expect(page).toHaveURL('/');
    console.log('✓ Login successful');

    // Step 2: Create a patient
    const timestamp = Date.now();
    const patientData = {
      name: `Integration Patient ${timestamp}`,
      yob: '1990',
      gender: 'Nam',
      address: '456 Integration Ave, Test City',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };

    await goToPatients(page);
    await createPatient(page, patientData);
    await searchPatient(page, patientData.name);
    await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
    console.log('✓ Patient created successfully');

    // Step 3: Create multiple tests
    const tests = [
      {
        name: `Integration Test 1 ${timestamp}`,
        unit: 'mmol/L',
        price: 50000,
        lowerBound: 3.5,
        upperBound: 6.5,
        normalValue: '3.5-6.5',
      },
      {
        name: `Integration Test 2 ${timestamp}`,
        unit: 'g/dL',
        price: 60000,
        lowerBound: 12.0,
        upperBound: 16.0,
        normalValue: '12.0-16.0',
      },
      {
        name: `Integration Test 3 ${timestamp}`,
        unit: 'mg/dL',
        price: 45000,
        lowerBound: 70,
        upperBound: 110,
        normalValue: '70-110',
      },
    ];

    await goToTests(page);
    for (const testData of tests) {
      await createTest(page, testData);
      await searchTest(page, testData.name);
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    }
    console.log('✓ All tests created successfully');

    // Step 4: Create a combo with the tests
    const comboData = {
      name: `Integration Combo ${timestamp}`,
      tests: tests.map(t => t.name),
    };

    await goToCombos(page);
    await createCombo(page, comboData);
    await searchCombo(page, comboData.name);
    await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
    console.log('✓ Combo created successfully');

    // Step 5: Create a record with the patient and combo
    const recordData = {
      patientName: patientData.name,
      comboName: comboData.name,
      testResults: [
        { testName: tests[0].name, value: '5.0' },  // Normal value
        { testName: tests[1].name, value: '14.5' }, // Normal value
        { testName: tests[2].name, value: '95' },   // Normal value
      ],
    };

    await createRecord(page, recordData);
    console.log('✓ Record created successfully');

    // Step 6: Verify the record appears in the records list
    await goToRecords(page);
    await page.getByPlaceholder('Tên bệnh nhân hoặc số điện thoại').fill(patientData.name);
    await page.waitForTimeout(600);
    
    const recordRow = page.locator('tr', { hasText: patientData.name }).first();
    const hasRecord = await recordRow.count() > 0;
    expect(hasRecord).toBeTruthy();
    console.log('✓ Record verified in records list');

    // Step 7: View record details
    if (hasRecord) {
      await recordRow.getByRole('link', { name: 'Xem' }).click();
      await page.waitForLoadState('networkidle');
      
      // Verify patient information
      await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
      console.log('✓ Record details page displayed correctly');
    }

    console.log('✓✓✓ Complete integration test passed successfully! ✓✓✓');
  });

  test.skip('should handle workflow with abnormal test results', async ({ page }) => {
    // Login
    await login(page);

    const timestamp = Date.now();

    // Create patient
    const patientData = {
      name: `Abnormal Test Patient ${timestamp}`,
      yob: '1988',
      gender: 'Nữ',
      address: '789 Abnormal St',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };

    await goToPatients(page);
    await createPatient(page, patientData);
    console.log('✓ Patient created');

    // Create a test with specific bounds
    const testData = {
      name: `Abnormal Test ${timestamp}`,
      unit: 'mmol/L',
      price: 50000,
      lowerBound: 3.0,
      upperBound: 6.0,
      normalValue: '3.0-6.0',
    };

    await goToTests(page);
    await createTest(page, testData);
    console.log('✓ Test created');

    // Create combo
    const comboData = {
      name: `Abnormal Combo ${timestamp}`,
      tests: [testData.name],
    };

    await goToCombos(page);
    await createCombo(page, comboData);
    console.log('✓ Combo created');

    // Create record with abnormal value (above upper bound)
    const recordData = {
      patientName: patientData.name,
      comboName: comboData.name,
      testResults: [
        { testName: testData.name, value: '8.5' }, // Abnormally high
      ],
    };

    await page.goto('/phieu-xet-nghiem/new');
    await page.waitForLoadState('networkidle');

    // Fill in patient using autocomplete
    const patientInput = page.getByRole('row', { name: 'Bệnh nhân' }).getByRole('textbox');
    await patientInput.fill(recordData.patientName);
    await page.waitForTimeout(600);
    
    const patientSuggestion = page.locator('.autocomplete-option', { hasText: recordData.patientName }).first();
    if (await patientSuggestion.count() > 0) {
      await patientSuggestion.click();
      await page.waitForTimeout(300);
    }

    // Fill in combo
    const comboInput = page.getByRole('row', { name: 'Tên gói xét nghiệm' }).getByRole('textbox');
    await comboInput.fill(recordData.comboName);
    await page.waitForTimeout(600);
    
    const comboSuggestion = page.locator('.autocomplete-option', { hasText: recordData.comboName }).first();
    if (await comboSuggestion.count() > 0) {
      await comboSuggestion.click();
      await page.waitForTimeout(500);
    }

    // Enter abnormal test value
    const testRow = page.locator('tr', { hasText: testData.name }).first();
    const testInput = testRow.getByRole('textbox').nth(0);
    
    if (await testInput.count() > 0) {
      await testInput.fill(recordData.testResults[0].value);
      console.log('✓ Abnormal value entered');
    }

    // Submit the form
    await page.getByRole('button', { name: 'Tạo (Ctrl + S)' }).click();
    await page.waitForTimeout(300);

    expect(testRow.getByRole('checkbox')).toBeChecked()

    console.log('✓✓✓ Abnormal value workflow test passed! ✓✓✓');
  });

  test('should navigate through all main pages', async ({ page }) => {
    await login(page);

    const pages = [
      { url: '/phieu-xet-nghiem', name: 'Records' },
      { url: '/danh-muc-benh-nhan', name: 'Patients' },
      { url: '/danh-muc-xet-nghiem', name: 'Tests' },
      { url: '/danh-muc-goi-xet-nghiem', name: 'Combos' },
      { url: '/danh-muc-bac-si', name: 'Doctors' },
    ];

    for (const pageInfo of pages) {
      await page.goto(pageInfo.url);
      await page.waitForLoadState('networkidle');
      
      // Verify page loaded by checking for main content
      const table = page.getByRole('table');
      await expect(table).toBeVisible();
      
      console.log(`✓ Navigated to ${pageInfo.name} (${pageInfo.url})`);
    }

    console.log('✓✓✓ All pages navigation test passed! ✓✓✓');
  });
});
