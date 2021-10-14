describe('2fa mandatory app', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'mandatory' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it('2fa flow mandatory with app', () => {
    cy.loginWithEmail();
    cy.get('.tfa-setup__qr-code');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.contains('Success');
    cy.screenshot();
  });
});
