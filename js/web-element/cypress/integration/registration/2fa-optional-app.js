describe('2fa optional registration', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'optional' });
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('ask for setup 2fa and setup app 2fa', () => {
    cy.registerWithEmail();
    cy.contains('Authenticator app');
    cy.screenshot();
    cy.get('button').contains('Setup').click();
    cy.get('button').contains('Continue');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.contains('Success');
    cy.screenshot();
  });
  it('ask for setup 2fa and skip', () => {
    cy.registerWithEmail();
    cy.contains('Setup next time');
    cy.screenshot();
    cy.contains('Setup next time').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
