describe('simple login', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.appSet({ tfa_status: 'disabled' });
    cy.visitLogin();
  });
  it('login by email', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Success');
  });
});
