describe('2fa mandatory app', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow disabled with app', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.get('.tfa-setup__qr-code');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.get('#tfaCode').click().type('0000');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
