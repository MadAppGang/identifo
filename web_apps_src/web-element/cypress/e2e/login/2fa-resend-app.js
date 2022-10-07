describe('2fa mandatory app', () => {
  before(() => {
    cy.createAppAndUser();
    cy.appSet({ tfa_status: 'mandatory' });
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app', tfa_resend_timeout: 5 });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it("2fa flow mandatory with app don't have resend", () => {
    cy.loginWithEmail();
    cy.get('.tfa-setup__qr-code');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.wait(1000);
    cy.contains('Resend code').should('not.exist');
    cy.wait(2000);
    cy.contains('Resend code').should('not.exist');
  });
});
