describe('simple login by email', () => {
  before(() => {
    cy.serverSetLoginOptions({});
    cy.appSet({ tfa_status: 'disabled' });
    cy.addTestUser();
  });
  it('login by email disabled', () => {
    cy.visitLogin();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Login').click();
    cy.screenshot();
    cy.contains('Application does not support login with email');
  });
  it('enable login by email', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.visitLogin();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
