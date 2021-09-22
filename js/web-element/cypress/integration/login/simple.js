describe('simple login', () => {
  beforeEach(() => {
    cy.disable2fa();
    cy.visitLogin();
  });
  it('login by email', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Success');
    cy.contains('test@test.com');
  });
});
