describe('2fa mandatory email', () => {
  before(() => {
    cy.createApp();
    cy.appSet({ tfa_status: 'mandatory' });
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('Login with app where 2fa disabled', () => {
    // register with 2fa
    cy.registerWithEmail();
    cy.get('button').contains('Continue');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
    cy.visitLogin();

    // login
    cy.appSet({ tfa_status: 'disabled' });
    cy.loginWithEmail();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
