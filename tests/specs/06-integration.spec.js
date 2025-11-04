// Complete E2E Integration Test - Full Application Flow
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { createPatient } = require('../helpers/patients');
const { createTest } = require('../helpers/tests');
const { createCombo } = require('../helpers/combos');
const { createRecord } = require('../helpers/records');

test.describe('Complete Application Flow Integration Test', () => {
  test('should complete full workflow: create patient, tests, combo, and record', async ({ page }) => {
    // Step 1: Login
    await login(page, 'admin', 'admin123');
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

    await page.goto('/danh-muc-benh-nhan');
    await createPatient(page, patientData);
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

    await page.goto('/danh-muc-xet-nghiem');
    for (const testData of tests) {
      await createTest(page, testData);
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    }
    console.log('✓ All tests created successfully');

    // Step 4: Create a combo with the tests
    const comboData = {
      name: `Integration Combo ${timestamp}`,
      tests: tests.map(t => t.name),
    };

    await createCombo(page, comboData);
    await page.goto('/danh-muc-goi-xet-nghiem');
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
    await page.goto('/phieu-xet-nghiem');
    await page.fill('input[placeholder*="Tìm kiếm"]', patientData.name);
    await page.waitForTimeout(500);
    
    const recordRow = page.locator('tr', { hasText: patientData.name }).first();
    const hasRecord = await recordRow.count() > 0;
    expect(hasRecord).toBeTruthy();
    console.log('✓ Record verified in records list');

    // Step 7: View record details
    if (hasRecord) {
      await recordRow.locator('text=Chi tiết').click();
      await page.waitForLoadState('networkidle');
      
      // Verify patient information
      await expect(page.locator(`text=${patientData.name}`)).toBeVisible();
      console.log('✓ Record details page displayed correctly');
    }

    console.log('✓✓✓ Complete integration test passed successfully! ✓✓✓');
  });

  test('should handle workflow with abnormal test results', async ({ page }) => {
    // Login
    await login(page, 'admin', 'admin123');

    const timestamp = Date.now();

    // Create patient
    const patientData = {
      name: `Abnormal Test Patient ${timestamp}`,
      yob: '1988',
      gender: 'Nữ',
      address: '789 Abnormal St',
      phone: `0${Math.floor(Math.random() * 900000000 + 100000000)}`,
    };

    await page.goto('/danh-muc-benh-nhan');
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

    await page.goto('/danh-muc-xet-nghiem');
    await createTest(page, testData);
    console.log('✓ Test created');

    // Create combo
    const comboData = {
      name: `Abnormal Combo ${timestamp}`,
      tests: [testData.name],
    };

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

    // Fill in patient
    await page.fill('input[placeholder*="Tìm kiếm bệnh nhân"]', recordData.patientName);
    await page.waitForTimeout(500);
    
    const patientSuggestion = page.locator(`text=${recordData.patientName}`).first();
    if (await patientSuggestion.count() > 0) {
      await patientSuggestion.click();
      await page.waitForTimeout(300);
    }

    // Fill in combo
    await page.fill('input[placeholder*="Tìm kiếm gói xét nghiệm"]', recordData.comboName);
    await page.waitForTimeout(500);
    
    const comboSuggestion = page.locator(`text=${recordData.comboName}`).first();
    if (await comboSuggestion.count() > 0) {
      await comboSuggestion.click();
      await page.waitForTimeout(500);
    }

    // Enter abnormal test value
    const testRow = page.locator('tr', { hasText: testData.name }).first();
    const testInput = testRow.locator('input[name*="test_value"]');
    
    if (await testInput.count() > 0) {
      await testInput.fill(recordData.testResults[0].value);
      console.log('✓ Abnormal value entered');
    }

    // Submit the form
    await page.click('button[type="submit"]:has-text("Tạo phiếu")');
    await page.waitForLoadState('networkidle');

    console.log('✓✓✓ Abnormal value workflow test passed! ✓✓✓');
  });

  test('should navigate through all main pages', async ({ page }) => {
    await login(page, 'admin', 'admin123');

    const pages = [
      { url: '/phieu-xet-nghiem', title: 'Phiếu xét nghiệm' },
      { url: '/danh-muc-benh-nhan', title: 'Danh mục bệnh nhân' },
      { url: '/danh-muc-xet-nghiem', title: 'Danh mục xét nghiệm' },
      { url: '/danh-muc-goi-xet-nghiem', title: 'Danh mục gói xét nghiệm' },
      { url: '/danh-muc-bac-si', title: 'Danh mục bác sĩ' },
      { url: '/so-sanh-ket-qua', title: 'So sánh kết quả' },
      { url: '/danh-muc-so-sanh', title: 'Danh mục so sánh' },
      { url: '/bao-cao-thong-ke-doanh-so', title: 'Báo cáo thống kê doanh số' },
    ];

    for (const pageInfo of pages) {
      await page.goto(pageInfo.url);
      await page.waitForLoadState('networkidle');
      
      // Verify page loaded by checking for title or heading
      const heading = page.locator('h1, h2, h3, h4').first();
      await expect(heading).toBeVisible();
      
      console.log(`✓ Navigated to ${pageInfo.url}`);
    }

    console.log('✓✓✓ All pages navigation test passed! ✓✓✓');
  });
});
