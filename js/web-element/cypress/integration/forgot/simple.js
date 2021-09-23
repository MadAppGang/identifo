describe('simple forgot', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({});
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
    cy.screenshot();
  });
});
