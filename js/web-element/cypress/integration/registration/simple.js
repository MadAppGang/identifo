describe('simple registration', () => {
  beforeEach(() => {
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('open signup and register test email', () => {
    cy.contains('Sign Up').click();
    cy.get('#login');
    cy.get('#password');
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Continue').click();
    cy.contains('Success');
    cy.contains('test@test.com');
  });
});
