describe('simple login by email', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('login by email', () => {
    cy.visitLogin();
    cy.loginWithEmail();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
