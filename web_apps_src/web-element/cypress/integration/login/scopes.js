describe('Check scopes pass', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('Login with empty scopes', () => {
    cy.appSet({ scopes: [] });
    cy.userSet({ scopes: [] });
    cy.visitLogin({ scopes: '' });
    cy.loginWithEmail();
    cy.screenshot();
  });
  it('Login with scope on user only', () => {
    cy.appSet({ scopes: [] });
    cy.userSet({ scopes: ['checkMe'] });
    cy.visitLogin({ scopes: 'checkMe' });
    cy.loginWithEmail();
    cy.contains('"scopes":["checkMe"]').click();
  });
  it('Login with scope on app only', () => {
    cy.appSet({ scopes: ['checkMe'] });
    cy.userSet({ scopes: [] });
    cy.visitLogin({ scopes: 'checkMe' });
    cy.loginWithEmail();
    cy.contains('Internal server error').click();
  });
  it('Login with scope on user and app', () => {
    cy.appSet({ scopes: ['checkMe'] });
    cy.userSet({ scopes: ['checkMe'] });
    cy.visitLogin({ scopes: 'checkMe' });
    cy.loginWithEmail();
    cy.contains('"scopes":["checkMe"]').click();
  });
});
