describe('login by phone', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: true, email: false, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('login by phone', () => {
    cy.userSet({ phone: '+0123456789' });
    cy.visitLogin();
    cy.loginWithPhone();
    cy.loginWithPhoneVerify();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
