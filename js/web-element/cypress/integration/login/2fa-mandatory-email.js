describe('2fa mandatory email', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'email' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow disabled with email', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Use email as 2fa');
    cy.screenshot();
    cy.get('button').contains('Setup email').click();
    cy.get('#tfaCode').click().type('0000');
    cy.screenshot();
    cy.contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
