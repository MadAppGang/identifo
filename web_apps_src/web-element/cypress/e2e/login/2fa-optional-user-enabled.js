describe('2fa mandatory email', () => {
  before(() => {
    cy.createApp();
    cy.appSet({ tfa_status: 'optional' });

    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('Ask 2fa with app where 2fa optional and user enabled', () => {
    // register with 2fa
    cy.registerWithEmail();
    cy.contains('Authenticator app');
    cy.screenshot();
    cy.get('button').contains('Setup').click();
    cy.get('button').contains('Continue');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
    cy.visitLogin();
    // login
    cy.loginWithEmail();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
