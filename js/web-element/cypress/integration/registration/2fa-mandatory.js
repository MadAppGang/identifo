describe('simple registration', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({});
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('open signup and register test email', () => {
    cy.screenshot();
    cy.contains('Sign Up').click();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.contains('Continue').click();
    cy.contains('Authenticator app');
    cy.screenshot();
  });
});
