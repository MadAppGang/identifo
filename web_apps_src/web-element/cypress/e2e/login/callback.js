describe('simple login by email', () => {
  before(() => {
    cy.serverSetLoginOptions({});
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('check error callback', () => {
    cy.visitLogin({ callbackUrl: 'https://google.com' });
    cy.contains('Please add callbackURL');
    cy.screenshot();
  });
  it('check callback works and have token', () => {
    cy.appSet({ redirect_urls: ['http://localhost:44000', 'https://google.com'], tfa_status: 'disabled' });
    cy.visitLogin({ callbackUrl: 'https://google.com' });
    cy.loginWithEmail();
    cy.location('host').should('eq', 'www.google.com');
    cy.location('search').should(search => {
      expect(search.match(/token=[\w\d]+/)).to.have.length(1);
    });
  });
  it('check callback works and have token and refresh-token when scopes offline', () => {
    cy.appSet({ redirect_urls: ['http://localhost:44000', 'https://google.com'], tfa_status: 'disabled' });
    cy.visitLogin({ callbackUrl: 'https://google.com', scopes: 'offline' });
    cy.loginWithEmail();
    cy.location('host').should('eq', 'www.google.com');
    cy.location('search').should(search => {
      expect(search.match(/token=[\w\d]+/)).to.have.length(1);
      expect(search.match(/refresh_token=[\w\d]+/)).to.have.length(1);
    });
  });
});
