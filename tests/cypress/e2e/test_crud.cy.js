describe('template spec', () => {
  it('passes', () => {
    cy.visit('http://127.0.0.1:7331/')
  })
})

describe('Test CRUD flow', () => {
  const testData = [
    { name: 'Glucose', unit: 'mmol/L', price: 1000, lower: 3.9, upper: 6.1, normal: '3.9-6.1' },
    { name: 'Cholesterol', unit: 'mmol/L', price: 1200, lower: 3.0, upper: 5.2, normal: '3.0-5.2' },
  ];

  it('should create 10 tests and verify they appear in the list', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-xet-nghiem');
    testData.forEach((test, idx) => {
      cy.contains('Thêm xét nghiệm').click();
      cy.get('input#test_name-input[name="test_name"]').type(test.name);
      cy.get('input#test_unit-input[name="test_unit"]').type(test.unit);
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