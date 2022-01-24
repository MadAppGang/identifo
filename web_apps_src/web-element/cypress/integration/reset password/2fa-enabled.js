describe('simple reset password with tfa email', () => {
  before(() => {
    cy.createAppAndUser();
    cy.appSet({ tfa_status: 'mandatory' });
    cy.userSet({ phone: '+1234567890', tfa_info: { is_enabled: true } });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it('forgot by email with 2fa email', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'email' });
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.screenshot();
    cy.contains('Send the link').click();
    cy.contains('We sent you an email with a link');
    cy.contains('Go back to login');
    cy.screenshot();
  });
  it('forgot by email with 2fa app', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.screenshot();
    cy.contains('Send the link').click();
    cy.verifyTfa();
    cy.contains('We sent you an email with a link');
    cy.contains('Go back to login');
    cy.screenshot();
  });
  it('forgot by email with 2fa sms', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.screenshot();
    cy.contains('Send the link').click();
    cy.verifyTfa();
    cy.contains('We sent you an email with a link');
    cy.contains('Go back to login');
    cy.screenshot();
  });
});
