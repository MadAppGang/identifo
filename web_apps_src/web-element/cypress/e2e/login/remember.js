describe('simple login by email', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({});
    cy.appSet({ tfa_status: 'disabled' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('login with remember', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.visitLogin();
    cy.loginWithEmail({ remember: true });
    cy.verifyRefreshSuccessToken();
    cy.screenshot();
  });
});
