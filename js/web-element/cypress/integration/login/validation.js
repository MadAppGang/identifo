describe('login errors', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: '' });
    cy.appSet({});
    cy.visitLogin();
  });
  it('login button disabled when empty user or password', () => {
    cy.contains('Login').should('be.disabled');
    cy.get('#login').click().type('test');
    cy.contains('Login').should('be.disabled');
    cy.get('#password').click().type('test');
    cy.contains('Login').should('be.not.disabled');
    cy.screenshot();
  });
  it('validate email', () => {
    cy.get('#login').click().type('fake user');
    cy.get('#password').click().type('wrong password');
    cy.contains('Login').click();
    cy.get('.error').should('have.text', 'Email address is not valid');
    cy.screenshot();
  });
  it('invalid user', () => {
    cy.get('#login').click().type('fake@test.com');
    cy.get('#password').click().type('wrong password');
    cy.contains('Login').click();
    cy.get('.error').should('have.text', 'Incorrect email or password');
    cy.screenshot();
  });
});
