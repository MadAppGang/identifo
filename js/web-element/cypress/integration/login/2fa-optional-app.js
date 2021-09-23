// TASK: When optional remember user to enable 2fa
describe('2fa optional', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'optional' });
    cy.addTestUser();
    cy.visitLogin();
  });
  it('ask for setup 2fa with skip avaible', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Skip');
  });
});
