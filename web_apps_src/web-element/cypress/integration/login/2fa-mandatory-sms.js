describe('2fa mandatory sms', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow mandatory with empty phone', () => {
    cy.addTestUser();
    cy.loginWithEmail();
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.verifyTfa();
    cy.contains('Success');
    cy.screenshot();
  });
  it('2fa flow mandatory with phone filled', () => {
    cy.addTestUser({ phone: '+0123456789' });
    cy.loginWithEmail();
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').should('have.value', '+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.verifyTfa();
    cy.contains('Success');
    cy.screenshot();
  });
});
