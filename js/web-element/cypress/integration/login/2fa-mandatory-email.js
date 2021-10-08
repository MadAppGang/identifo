describe('2fa mandatory email', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'email' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow mandatory with email', () => {
    cy.loginWithEmail();
    cy.contains('Use email as 2fa');
    cy.screenshot();
    cy.get('button').contains('Setup email').click();
    cy.verifyTfa();
    cy.contains('Success');
    cy.screenshot();
  });
});
