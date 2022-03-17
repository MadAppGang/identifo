describe('2fa mandatory email', () => {
  before(() => {
    cy.createApp();
    cy.createUser({ tfa_info: { is_enabled: true } });
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({ tfa_status: 'mandatory' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('forgot by email with 2fa sms', () => {
    cy.visitLogin();
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.screenshot();
    cy.contains('Send the link').click();
    cy.contains('We sent you an email with a link');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.getResetTokenURL().then(url => {
      cy.visit(url);
    });
    cy.contains('Set up a new password');
    cy.get('#password').click().type('NewPassword');
    cy.screenshot();
    cy.contains('Save password').click().type('NewPassword');

    cy.loginWithEmail({ password: 'NewPassword' });
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('[placeholder=Phone]').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
