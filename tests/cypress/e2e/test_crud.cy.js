describe('template spec', () => {
  it('passes', () => {
    cy.visit('http://127.0.0.1:7331/')
  })
})

describe('Test CRUD flow', () => {
  const testData = [
    { name: 'Glucose', unit: 'mmol/L', price: 1000, lower: 3.9, upper: 6.1, normal: '3.9-6.1' },
    { name: 'Cholesterol', unit: 'mmol/L', price: 1200, lower: 3.0, upper: 5.2, normal: '3.0-5.2' },
    { name: 'Triglyceride', unit: 'mmol/L', price: 1100, lower: 0.4, upper: 1.7, normal: '0.4-1.7' },
    { name: 'Urea', unit: 'mmol/L', price: 900, lower: 2.5, upper: 7.5, normal: '2.5-7.5' },
    { name: 'Creatinine', unit: 'umol/L', price: 950, lower: 62, upper: 120, normal: '62-120' },
    { name: 'ALT (GPT)', unit: 'U/L', price: 1050, lower: 0, upper: 40, normal: '0-40' },
    { name: 'AST (GOT)', unit: 'U/L', price: 1050, lower: 0, upper: 40, normal: '0-40' },
    { name: 'Bilirubin', unit: 'umol/L', price: 1150, lower: 3.4, upper: 20.5, normal: '3.4-20.5' },
    { name: 'Protein', unit: 'g/L', price: 1000, lower: 60, upper: 80, normal: '60-80' },
    { name: 'Albumin', unit: 'g/L', price: 1000, lower: 35, upper: 50, normal: '35-50' }
  ];

  it('should create 10 tests and verify they appear in the list', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-xet-nghiem');
    testData.forEach((test, idx) => {
      cy.contains('Thêm xét nghiệm').click();
      cy.get('input#test_name[name="test_name"]').type(test.name);
      cy.get('input[name="test_unit"]').type(test.unit);
      cy.get('input[name="test_price"]').type(`${test.price}`);
      cy.get('input[name="test_lower_bound"]').type(`${test.lower}`);
      cy.get('input[name="test_upper_bound"]').type(`${test.upper}`);
      cy.get('input[name="test_normal_value"]').type(test.normal);
      cy.get('button[type="submit"]').contains('Thêm xét nghiệm').click();
      cy.contains(test.name).should('exist');
    });
  });

  it('should clean up all test data by searching and deleting each test', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-xet-nghiem');
    testData.forEach((test) => {
      cy.get('#test-search').clear().type(test.name);
      cy.get('#test-search-form').submit();
      cy.contains(test.name).parents('tr').within(() => {
        cy.contains('Xoá').click();
      });
      cy.contains(test.name).should('not.exist');
    });
  });

});