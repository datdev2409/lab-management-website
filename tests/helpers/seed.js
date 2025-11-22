// Helper functions for seeding test data
const { login } = require('./auth');
const { goToTests, createTest } = require('./tests');
const { goToCombos, createCombo } = require('./combos');

/**
 * Seed basic test data (tests and combos) for use in record tests
 * @param {import('@playwright/test').Page} page
 * @returns {Promise<{tests: Array, combos: Array}>} Created test and combo data
 */
async function seedBasicTestData(page) {
  await login(page);

  const timestamp = Date.now();
  
  // Create standard tests
  const tests = [
    {
      name: `Glucose ${timestamp}`,
      unit: 'mmol/L',
      price: 50000,
      lowerBound: 3.9,
      upperBound: 6.1,
      normalValue: '3.9-6.1',
    },
    {
      name: `Hemoglobin ${timestamp}`,
      unit: 'g/dL',
      price: 40000,
      lowerBound: 12.0,
      upperBound: 16.0,
      normalValue: '12.0-16.0',
    },
    {
      name: `WBC Count ${timestamp}`,
      unit: '10^9/L',
      price: 35000,
      lowerBound: 4.0,
      upperBound: 11.0,
      normalValue: '4.0-11.0',
    },
    {
      name: `Cholesterol ${timestamp}`,
      unit: 'mmol/L',
      price: 60000,
      lowerBound: 3.0,
      upperBound: 5.2,
      normalValue: '3.0-5.2',
    },
  ];

  await goToTests(page);
  for (const testData of tests) {
    await createTest(page, testData);
  }

  // Create standard combos
  const combos = [
    {
      name: `Basic Health Check ${timestamp}`,
      tests: [tests[0].name, tests[1].name],
    },
    {
      name: `Complete Blood Count ${timestamp}`,
      tests: [tests[1].name, tests[2].name],
    },
    {
      name: `Full Panel ${timestamp}`,
      tests: tests.map(t => t.name),
    },
  ];

  await goToCombos(page);
  for (const comboData of combos) {
    await createCombo(page, comboData);
  }

  return { tests, combos };
}

module.exports = {
  seedBasicTestData,
};
