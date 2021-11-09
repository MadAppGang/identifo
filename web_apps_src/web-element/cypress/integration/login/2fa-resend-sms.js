describe('2fa mandatory app', () => {
  before(() => {
    cy.createAppAndUser();
    cy.appSet({ tfa_status: 'mandatory' });
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms', tfa_resend_timeout: 5 });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it('2fa flow mandatory with sms have resend', () => {
    cy.loginWithEmail();
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.wait(1000);
    cy.contains('Resend code').should('not.exist');
    cy.wait(2000);
    cy.contains('Resend code').should('exist');
    cy.contains('Resend code').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
