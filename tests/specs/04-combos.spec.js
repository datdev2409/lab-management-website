// Combo (Test Package) Management E2E tests
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToCombos, createCombo, searchCombo, deleteCombo } = require('../helpers/combos');
const { goToTests, createTest } = require('../helpers/tests');

test.describe('Combo Management Flow', () => {
  let testNames = [];

  test.beforeAll(async ({ browser }) => {
    // Create some tests to use in combos
    const context = await browser.newContext();
    const page = await context.newPage();
    
    await login(page);
    await goToTests(page);
    
    const tests = [
      {
        name: `Combo Test A ${Date.now()}`,
        unit: 'mmol/L',
        price: 50000,
        lowerBound: 3.0,
        upperBound: 6.0,
        normalValue: '3.0-6.0',
      },
      {
        name: `Combo Test B ${Date.now()}`,
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

  test.beforeEach(async ({ page }) => {
    // Login before each test
    await login(page);
  });

  test('should display combo management page', async ({ page }) => {
    await goToCombos(page);
    
    await expect(page.getByRole('heading', { name: 'Danh mục gói xét nghiệm' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Tạo gói xét nghiệm mới' })).toBeVisible();
  });

  test('should create a new combo with multiple tests', async ({ page }) => {
    await goToCombos(page);
    
    const comboData = {
      name: `Basic Health Check ${Date.now()}`,
      tests: testNames, // Use the tests created in beforeAll
    };
    
    await createCombo(page, comboData);
    
    // Navigate back to combo list
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
      tests: [testNames[0]], // Use at least one test
    };
    
    await createCombo(page, comboData);
    await goToCombos(page);
    
    // Search for the combo
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
    
    // Click to view details
    await searchCombo(page, comboData.name);
    const row = page.locator('tr', { hasText: comboData.name }).first();
    await row.locator('text=Chi tiết').click();
    
    await page.waitForLoadState('networkidle');
    
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
    
    // Verify combo exists
    await searchCombo(page, comboData.name);
    await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
    
    // Delete the combo
    await deleteCombo(page, comboData.name);
    
    // Verify combo is removed
    await searchCombo(page, comboData.name);
    await expect(page.locator(`text=${comboData.name}`)).not.toBeVisible();
  });

  test('should create multiple combos with different test combinations', async ({ page }) => {
    await goToCombos(page);
    
    const combos = [
      {
        name: `Combo Alpha ${Date.now()}`,
        tests: [testNames[0]],
      },
      {
        name: `Combo Beta ${Date.now()}`,
        tests: [testNames[1]],
      },
      {
        name: `Combo Full ${Date.now()}`,
        tests: testNames,
      },
    ];
    
    for (const comboData of combos) {
      await createCombo(page, comboData);
      await goToCombos(page);
      await searchCombo(page, comboData.name);
      await expect(page.locator(`text=${comboData.name}`)).toBeVisible();
    }
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
    await page.waitForLoadState('networkidle');
    
    // Click back button
    await page.getByRole('button', { name: 'Trở lại' }).click();
    await page.waitForTimeout(500);
    
    // Verify we're back on the combo list page
    await expect(page).toHaveURL('/danh-muc-goi-xet-nghiem');
  });
});
