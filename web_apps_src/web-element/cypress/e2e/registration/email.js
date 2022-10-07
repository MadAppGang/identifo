describe('simple registration', () => {
  before(() => {
    cy.createAppAndUser(false);
    cy.serverSetLoginOptions({});
    cy.appSet({ tfa_status: 'disabled' });
    cy.deleteTestUser();
    cy.visitLogin();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('open signup and register test email', () => {
    cy.screenshot();
    cy.contains('Sign Up').click();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.contains('Continue').click();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
