describe('simple reset password without tfa', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
    cy.addTestUser();
    cy.visitLogin();
  });
  it('have back button', () => {
    cy.contains('Forgot password').click();
    cy.contains('Go back to login');
    cy.screenshot();
  });
  it('forgot by email', () => {
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.screenshot();
    cy.contains('Send the link').click();
    cy.contains('We sent you an email with a link to create a new password');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.getResetTokenURL().then(url => {
      cy.visit(url);
    });
    cy.contains('Set up a new password');
    cy.get('#password').click().type('NewPassword');
    cy.screenshot();
    cy.contains('Save password').click().type('NewPassword');

    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('NewPassword');
    cy.contains('Login').click();
  });
});
