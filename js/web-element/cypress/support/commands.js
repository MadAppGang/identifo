import user from '../fixtures/user.json';
import app from '../fixtures/app.json';
import server from '../fixtures/server.json';
const adminUrl = `${Cypress.config('serverUrl')}/admin`;
const login = () => {
  return fetch(`${adminUrl}/login`, {
    body: '{"email":"admin@admin.com","password":"password"}',
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
};
const getTestUser = async () => {
  const testUser = await fetch(`${adminUrl}/users?search=test`, {
    body: null,
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => result.users[0]);
  if (testUser) {
    return testUser;
  }
};
const deleteTestUser = async () => {
  const testUser = await fetch(`${adminUrl}/users?search=test`, {
    body: null,
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => result.users[0]);
  if (testUser) {
    return fetch(`${adminUrl}/users/${testUser.id}`, {
      body: null,
      method: 'DELETE',
      mode: 'cors',
      credentials: 'include',
    });
  }
};
Cypress.Commands.add('appSet', async data => {
  await login();
  await fetch(`${adminUrl}/apps/${Cypress.config('appId')}`, {
    body: JSON.stringify({
      ...app,
      ...data,
    }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  });
});
Cypress.Commands.add('serverSetLoginOptions', async data => {
  await login();
  await fetch(`${adminUrl}/settings`, {
    body: JSON.stringify({
      login: { ...server.login, ...data },
    }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  });
});

Cypress.Commands.add('deleteTestUser', async () => {
  await login();
  await deleteTestUser();
});

Cypress.Commands.add('addTestUser', async data => {
  await login();
  await deleteTestUser();
  const testUser = await fetch(`${adminUrl}/users/`, {
    body: JSON.stringify({ ...user, ...data }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
});

Cypress.Commands.add('getResetTokenURL', async () => {
  const user = await getTestUser();
  if (user) {
    const resetTokenData = await fetch(`${adminUrl}/users/generate_new_reset_token`, {
      body: JSON.stringify({ user_id: user.id, app_id: Cypress.config('appId') }),
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
    }).then(r => r.json());
    return `${Cypress.config('baseUrl')}/password/reset/?appId=${Cypress.config('appId')}&url=${Cypress.config('serverUrl')}&token=${resetTokenData.Token}`;
  }
  throw new Error("Can't open reset token");
});

Cypress.Commands.add('visitLogin', (orig, url, options) => {
  return cy.visit(`${Cypress.config('baseUrl')}/login/?appId=${Cypress.config('appId')}&url=${Cypress.config('serverUrl')}`);
});
Cypress.Commands.add('loginWithEmail', (email = 'test@test.com', password = 'Password') => {
  cy.get('[placeholder=Email]').click().type(email);
  cy.get('[placeholder=Password]').click().type(password);
  cy.screenshot();
  cy.get('button').contains('Login').click();
});
Cypress.Commands.add('verifyTfa', (code = '0000') => {
  cy.get('[placeholder="Verify code"]').click().type('0000');
  cy.contains('Go back to login');
  cy.screenshot();
  cy.contains('Confirm').click();
});
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
