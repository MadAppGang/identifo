describe('2fa mandatory sms', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.serverSetLoginOptions({ login_with: { username: false, phone: false, email: true, federated: false }, tfa_type: 'sms' });
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('2fa flow with empty phone', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.screenshot();
    cy.contains('Login').click();
    cy.contains('Use phone as 2fa');
    cy.get('#phone').click().type('+0123456789');
    cy.screenshot();
    cy.contains('Setup phone').click();
    cy.get('#tfaCode').click().type('0000');
    cy.screenshot();
    cy.contains('Confirm').click();
    cy.contains('Success');
    cy.screenshot();
  });
});
