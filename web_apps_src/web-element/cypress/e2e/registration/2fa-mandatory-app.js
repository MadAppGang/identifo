describe('2fa mandatory registration', () => {
  before(() => {
    cy.createApp();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'mandatory' });
  });
  beforeEach(() => {
    cy.deleteTestUser();
    cy.visitLogin();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('open signup and register test email', () => {
    cy.registerWithEmail();
    cy.get('.tfa-setup__qr-code');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
