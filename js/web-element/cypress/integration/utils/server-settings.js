describe('server settings sync', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({});
  });
  it("check that federated login is hidden when isn't set", () => {
    cy.appSet({ federated_login_settings: {} });
    cy.visitLogin();
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
});
