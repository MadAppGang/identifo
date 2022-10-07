describe('server settings sync', () => {
  before(() => {
    cy.createAppAndUser();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
  });
  it("check that federated login is hidden when isn't set", () => {
    cy.appSet({ federated_login_settings: {} });
    cy.visitLogin();
    cy.contains('login with').should('not.exist');
    cy.screenshot();
  });
  it('check that federated login is shown when set apple', () => {
    cy.appSet({ federated_login_settings: { apple: { params: {}, scopes: [] } } });
    cy.visitLogin();
    cy.get('.social-buttons__apple');
    cy.screenshot();
  });
  it('check that federated login is shown when set google', () => {
    cy.appSet({ federated_login_settings: { google: { params: {}, scopes: [] } } });
    cy.visitLogin();
    cy.get('.social-buttons__google');
    cy.screenshot();
  });
  it('check that federated login is shown when set facebook', () => {
    cy.appSet({ federated_login_settings: { facebook: { params: {}, scopes: [] } } });
    cy.visitLogin();
    cy.get('.social-buttons__facebook');
    cy.screenshot();
  });
  it('check that sign is shown when not forbidden', () => {
    cy.appSet({ registration_forbidden: false });
    cy.visitLogin();
    cy.get('.login-form__register-link').should('exist');
    cy.screenshot();
  });
  it('check that sign up not shown when forbidden', () => {
    cy.appSet({ registration_forbidden: true });
    cy.visitLogin();
    cy.get('.login-form__register-link').should('not.exist');
    cy.screenshot();
  });
  it('check when login with phone is shown when login with phone and email enabled', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: true, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ registration_forbidden: true });
    cy.visitLogin();
    cy.get('[placeholder="Phone number"]').should('exist');
    cy.contains('login with password').click();
    cy.get('[placeholder="Password"]').should('exist');
    cy.screenshot();
  });
  it('check when OTP form is shown when login with is phone', () => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: true, email: false, federated: false }, tfa_type: 'app' });
    cy.appSet({ registration_forbidden: true });
    cy.visitLogin();
    cy.get('[placeholder="Phone number"]').should('exist');
    cy.contains('login with').should('not.exist');
    cy.screenshot();
  });
});
