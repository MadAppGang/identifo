describe('simple registration', () => {
  before(() => {
    cy.createAppAndUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'app' });
    cy.appSet({ tfa_status: 'disabled' });
  });
  beforeEach(() => {
    cy.deleteTestUser();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  it('open signup and register with wrong token', () => {
    cy.visitLogin({ token: 'test-token' });
    cy.screenshot();
    cy.contains('Sign Up').click();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.contains('Continue').click();
    cy.contains('token contains an invalid number of segments');
    cy.screenshot();
  });
  it('open signup and register with right token', () => {
    cy.addInvite('test@test.com', 'myrole').then(url => {
      cy.visit(url);
    });
    cy.get('#password').click().type('Password');
    cy.contains('Continue').click();
    cy.screenshot();
    cy.verifySuccessToken();
    cy.getUserData().then(user => {
      expect(user.access_role).equal('myrole');
    });
    cy.screenshot();
  });
});
