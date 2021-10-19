// TASK: When optional remember user to enable 2fa
describe('2fa optional', () => {
  before(() => {
    cy.createAppAndUser();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'email' });
    cy.appSet({ tfa_status: 'optional' });
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
  it('ask for setup 2fa and setup email 2fa', () => {
    cy.loginWithEmail();
    cy.contains('Use email as 2fa');
    cy.screenshot();
    cy.get('button').contains('Setup').click();
    cy.get('[placeholder=Email]').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup email').click();
    cy.verifyTfa();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});
