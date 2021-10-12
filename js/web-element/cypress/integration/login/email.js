describe('simple login by email', () => {
  before(() => {
    cy.serverSetLoginOptions({});
    cy.appSet({ tfa_status: 'disabled' });
    cy.addTestUser();
  });
  it('login by email disabled', () => {
    cy.visitLogin();
    cy.loginWithEmail();
    cy.contains('Login with username is not supported by app');
    cy.screenshot();
  });
  it('enable login by email', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.visitLogin();
    cy.loginWithEmail();
    cy.contains('Success');
    cy.screenshot();
  });
});
