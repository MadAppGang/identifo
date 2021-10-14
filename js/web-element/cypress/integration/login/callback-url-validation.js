describe('login errors', () => {
  before(() => {
    cy.createAppAndUser();
  });
  after(() => {
    cy.deleteAppAndUser();
  });
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: '' });
    cy.appSet({});
  });
  it('empty callback', () => {
    cy.visitLogin();
    cy.get('button').contains('Login');
  });
  it('valid callback', () => {
    cy.visitLogin({ callbackUrl: 'http://localhost:44000' });
    cy.get('button').contains('Login');
  });
  it('invalid callback', () => {
    cy.visitLogin({ callbackUrl: 'https://google.com' });
    cy.contains('Callback url is invalid');
  });
});
