import { app } from './app';
const login = () => {
  return fetch('https://identifo.organaza.ru/admin/login', {
    body: '{"email":"admin@admin.com","password":"password"}',
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
};
Cypress.Commands.add('disable2fa', async () => {
  await login();
  fetch(`https://identifo.organaza.ru/admin/apps/${Cypress.config('appId')}`, {
    body: JSON.stringify({
      ...app,
      tfa_status: 'disabled',
    }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  });
});

Cypress.Commands.add('deleteTestUser', async () => {
  await login();
  const testUser = await fetch('https://identifo.organaza.ru/admin/users?search=test', {
    body: null,
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => result.users[0]);
  if (testUser) {
    await fetch(`https://identifo.organaza.ru/admin/users/${testUser.id}`, {
      body: null,
      method: 'DELETE',
      mode: 'cors',
      credentials: 'include',
    });
  }
});
Cypress.Commands.add('visitLogin', (orig, url, options) => {
  return cy.visit(`${Cypress.config('baseUrl')}/?appId=${Cypress.config('appId')}&url=${Cypress.config('serverUrl')}&debug=true`);
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
