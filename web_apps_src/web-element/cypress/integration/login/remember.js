describe('simple login by email', () => {
  before(() => {
    cy.createAppAndUser();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  before(() => {
    cy.serverSetLoginOptions({});
    cy.appSet({ tfa_status: 'disabled' });
  });
  it('login with remember', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.visitLogin();
    cy.loginWithEmail(undefined, undefined, true);
    cy.verifyRefreshSuccessToken();
    cy.screenshot();
  });
});
