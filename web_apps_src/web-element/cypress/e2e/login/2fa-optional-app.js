// TASK: When optional remember user to enable 2fa
describe('2fa optional', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'optional' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.visitLogin();
  });
  it('ask for setup 2fa and skip', () => {
    cy.loginWithEmail();
    cy.contains('Setup next time');
    cy.screenshot();
    cy.contains('Setup next time').click();
    cy.verifySuccessToken();
    cy.screenshot();
  });
  it('ask for setup 2fa and setup app 2fa', () => {
    cy.loginWithEmail();
    cy.contains('Authenticator app');
    cy.screenshot();
    cy.get('button').contains('Setup').click();
    cy.get('button').contains('Continue');
    cy.screenshot();
    cy.get('button').contains('Continue').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
