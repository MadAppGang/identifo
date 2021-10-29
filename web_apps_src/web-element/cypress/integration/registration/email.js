describe('simple registration', () => {
  before(() => {
    cy.createAppAndUser(false);
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.serverSetLoginOptions({});
    cy.appSet({});
    cy.deleteTestUser();
    cy.visitLogin();
  });
  it('open signup and register test email', () => {
    cy.screenshot();
    cy.contains('Sign Up').click();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.contains('Continue').click();
    cy.verifySuccessToken();
    cy.screenshot();
  });
});