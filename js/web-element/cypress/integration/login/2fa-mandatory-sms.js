describe('2fa mandatory sms', () => {
  beforeEach(() => {
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow disabled with empty phone', () => {
    cy.addTestUser();
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.get('button').contains('Login').click();
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('#phone').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.get('#tfaCode').click().type('0000');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
  it('2fa flow disabled with phone filled', () => {
    cy.addTestUser({ phone: '+0123456789' });
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.get('button').contains('Login').click();
    cy.contains('Use phone as 2fa');
    cy.contains('Go back to login');
    cy.get('#phone').should('have.value', '+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.get('#tfaCode').click().type('0000');
    cy.contains('Go back to login');
    cy.screenshot();
    cy.get('button').contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
