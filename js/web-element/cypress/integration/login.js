const appId="c3vqpr2se6nhtg9v1nu0";
Cypress.on('uncaught:exception', (err, runnable) => {
    return false
})
describe('identifo-form', () => {
    it('display login form', () => {
        cy.visit('/');
        cy.wait(1000);
        cy.screenshot('login');
    })
    it('validate email', () => {
        cy.visit('/');
        cy.get('#floatingInput').click().type('fake user');
        cy.get('#floatingPassword').click().type('wrong password');
        cy.contains('Login').click();
        cy.get('.error').should('have.text', 'Email address is not valid.');
        cy.screenshot('login.email_error');
    })
    it('invalid user', () => {
        cy.visit('/');
        cy.get('#floatingInput').click().type('test@test.com');
        cy.get('#floatingPassword').click().type('wrong password');
        cy.contains('Login').click();
        cy.get('.error').should('have.text', 'User not found.');
        cy.screenshot('login.user_not_found');
    })
})