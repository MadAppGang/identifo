describe('simple login', () => {
  beforeEach(() => {
    cy.addTestUser();
    cy.appSet({ tfa_status: 'mandatory' });
    cy.visitLogin();
  });
  it('ask for 2fa google after login', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Login').click();
    cy.contains('Authenticator app');
    cy.screenshot();
  });
  it('check that google auth display qr code', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Login').click();
    cy.contains('Authenticator app');
    cy.contains('Setup').click();
    cy.contains('QR-code');
    cy.screenshot();
  });
  it('enter invalid code and show error', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Login').click();
    cy.contains('Authenticator app');
    cy.contains('Setup').click();
    cy.get('#tfaCode').click().type('1234');
    cy.contains('Confirm').click();
    cy.contains('Invalid two-factor authentication code');
    cy.screenshot();
  });
  it('enter valid code and press enter', () => {
    cy.get('#login').click().type('test@test.com');
    cy.get('#password').click().type('Password');
    cy.contains('Login').click();
    cy.contains('Authenticator app');
    cy.contains('Setup').click();
    cy.get('#tfaCode').click().type('0000{enter}');
    cy.contains('Success');
  });
});
