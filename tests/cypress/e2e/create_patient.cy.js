describe('Create patient flow', () => {
  it('should create a new patient and show in the list', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-benh-nhan');
    cy.contains('Thêm bệnh nhân').click();
    cy.get('#user_create_form').within(() => {
      cy.get('input[name="patient_name"]').type('Nguyen Van Test');
      cy.get('input[name="patient_yob"]').type('1990');
      cy.get('input[name="patient_gender"][value="Nam"]').check();
      cy.get('input[name="patient_address"]').type('123 Test St');
      cy.get('input[name="patient_phone"]').type('0123456789');
      cy.root().submit();
    });
    cy.visit('http://127.0.0.1:7331/danh-muc-benh-nhan');
    cy.contains('Nguyen Van Test').should('exist');
    cy.contains('123 Test St').should('exist');
    cy.contains('0123456789').should('exist');
  });

  it('should search for patient and verify results', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-benh-nhan');
    // Search for a non-existent patient
    cy.get('#patient-search').clear().type('abcdefghd');
    cy.get('#patient-search-form').submit();
    cy.contains('abcdefghd').should('not.exist');

    // Search for the newly created patient
    cy.get('#patient-search').clear().type('Nguyen Van Test');
    cy.get('#patient-search-form').submit();
    cy.contains('Nguyen Van Test').should('exist');
  });

  it('should edit patient details', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-benh-nhan');
    // Search for the patient first
    cy.get('#patient-search').clear().type('Nguyen Van Test');
    cy.get('#patient-search-form').submit();
    cy.contains('Nguyen Van Test').should('exist');
    
    cy.contains('Nguyen Van Test').parents('tr').within(() => {
      cy.contains('Sửa').click();
      cy.get('input[name="patient_name"]').clear().type('Nguyen Van Updated');
      cy.get('input[name="patient_address"]').clear().type('456 Updated St');
      cy.get('input[name="patient_phone"]').clear().type('0987654321');
      cy.contains('Lưu').click();
    });
    cy.on('window:confirm', () => true);
    cy.contains('Nguyen Van Updated').should('exist');
    cy.contains('456 Updated St').should('exist');
    cy.contains('0987654321').should('exist');
  });

  it('should delete patient', () => {
    cy.visit('http://127.0.0.1:7331/danh-muc-benh-nhan');
    // Search for the updated patient first
    cy.get('#patient-search').clear().type('Nguyen Van Updated');
    cy.get('#patient-search-form').submit();
    cy.contains('Nguyen Van Updated').should('exist');
    
    cy.contains('Nguyen Van Updated').parents('tr').within(() => {
      cy.contains('Xoá').click();
    });
    cy.window().then((win) => {
      cy.stub(win, 'confirm').returns(true);
    });
    cy.contains('Nguyen Van Updated').should('not.exist');
  });
});