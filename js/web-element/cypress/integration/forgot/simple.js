describe('simple forbidden', () => {
  beforeEach(() => {
    cy.disable2fa();
    cy.visitLogin();
  });
  it('forbidden by email', () => {
    cy.contains('Forgot password').click();
    cy.get('#email').click().type('test@test.com');
    cy.contains('Send the link').click();
    cy.contains('We sent you an email with a link to create a new password');
    cy.screenshot();
  });
});
