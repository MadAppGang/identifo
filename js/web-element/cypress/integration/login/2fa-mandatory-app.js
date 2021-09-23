describe('2fa mandatory app', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow with app', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Authenticator app');
    cy.screenshot();
    cy.contains('Setup').click();
    cy.contains('QR-code');
    cy.get('#tfaCode').click().type('0000');
    cy.screenshot();
    cy.contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
