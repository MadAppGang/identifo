import { app } from './app';
let testUser = {};
const login = () => {
  return fetch('https://identifo.organaza.ru/admin/login', {
    body: '{"email":"admin@admin.com","password":"password"}',
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
};
const deleteTestUser = async () => {
  const testUser = await fetch('https://identifo.organaza.ru/admin/users?search=test', {
    body: null,
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => result.users[0]);
  if (testUser) {
    return fetch(`https://identifo.organaza.ru/admin/users/${testUser.id}`, {
      body: null,
      method: 'DELETE',
      mode: 'cors',
      credentials: 'include',
    });
  }
};
Cypress.Commands.add('appSet', async data => {
  await login();
  await fetch(`https://identifo.organaza.ru/admin/apps/${Cypress.config('appId')}`, {
    body: JSON.stringify({
      ...app,
      ...data,
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

Cypress.Commands.add('addTestUser', async (tfa = false) => {
  await login();
  await deleteTestUser();
  testUser = await fetch(`https://identifo.organaza.ru/admin/users/`, {
    body: JSON.stringify({
      username: 'test@test.com',
      password: 'Password',
      access_role: '',
    }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  }).then(r => r.json());
  testUser = await fetch(`https://identifo.organaza.ru/admin/users/${testUser.id}`, {
    body: JSON.stringify({
      id: testUser.id,
      username: 'test@test.com',
      email: 'test@test.com',
      phone: '+123456789',
      active: true,
      tfa_info: {
        is_enabled: tfa,
      },
    }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  }).then(r => r.json());
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
