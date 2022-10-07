describe('2fa mandatory sms', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({ tfa_status: 'mandatory' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it('2fa flow mandatory with empty phone', () => {
    cy.userSet({});
    cy.loginWithEmail();
    cy.contains('Your phone will be used for 2-step verification');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
  it('2fa flow mandatory with phone filled', () => {
    cy.userSet({ phone: '+0123456789' });
    cy.loginWithEmail();
    cy.contains('Your phone will be used for 2-step verification');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').should('have.value', '+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
