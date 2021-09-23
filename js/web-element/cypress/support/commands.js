import user from '../fixtures/user.json';
import app from '../fixtures/app.json';
import server from '../fixtures/server.json';
let testUser = {};
const adminUrl = `${Cypress.config('serverUrl')}/admin`;
const login = () => {
  return fetch(`${adminUrl}/login`, {
    body: '{"email":"admin@admin.com","password":"password"}',
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
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
  testUser = await fetch(`${adminUrl}/users/`, {
    body: JSON.stringify({ ...user, ...data }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  }); //.then(r => r.json());
  // testUser = await fetch(`${adminUrl}/users/${testUser.id}`, {
  //   body: JSON.stringify({
  //     ...user,
  //     tfa_info: {
  //       is_enabled: tfa,
  //     },
  //   }),
  //   method: 'PUT',
  //   mode: 'cors',
  //   credentials: 'include',
  // }).then(r => r.json());
});

Cypress.Commands.add('visitLogin', (orig, url, options) => {
  return cy.visit(`${Cypress.config('baseUrl')}/?appId=${Cypress.config('appId')}&url=${Cypress.config('serverUrl')}`);
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
