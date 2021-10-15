import user from '../fixtures/user.json';
import app from '../fixtures/app.json';
import server from '../fixtures/server.json';
const adminUrl = `${Cypress.config('serverUrl')}/admin`;
let appId = '';
let userId = '';
const login = () => {
  return fetch(`${adminUrl}/login`, {
    body: '{"email":"admin@admin.com","password":"password"}',
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  });
};
const deleteTestUserBySearch = async () => {
  const testUser = await fetch(`${adminUrl}/users?search=test@test.com`, {
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
const deleteTestUser = async () => {
  if (!userId) { return }
  return fetch(`${adminUrl}/users/${userId}`, {
    body: null,
    method: 'DELETE',
    mode: 'cors',
    credentials: 'include',
  });
};
const deleteTestApp = async () => {
  if (!appId) { return }

  return fetch(`${adminUrl}/apps/${appId}`, {
    body: null,
    method: 'DELETE',
    mode: 'cors',
    credentials: 'include',
  });
};
// Init test user and app
Cypress.Commands.add('createAppAndUser', async (createUser = true) => {
  await login();
  await deleteTestApp();
  await deleteTestUser();
  await deleteTestUserBySearch();
  await fetch(`${adminUrl}/apps`, {
    body: JSON.stringify({
      ...app,
    }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => (appId = result.id));
  if (createUser) {
    await fetch(`${adminUrl}/users/`, {
      body: JSON.stringify({ ...user }),
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
    })
      .then(r => r.json())
      .then(result => (userId = result.id));
  }
});
// Cleanup test user and app
Cypress.Commands.add('deleteAppAndUser', async data => {
  await login();
  await deleteTestApp();
  await deleteTestUser();
});
// Change app settings
Cypress.Commands.add('appSet', async data => {
  await login();
  await fetch(`${adminUrl}/apps/${appId}`, {
    body: JSON.stringify({
      ...app,
      ...data,
    }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  });
});
// Change server login options
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

// Delete test user using search
// Using on registration tests
Cypress.Commands.add('deleteTestUser', async () => {
  await login();
  await deleteTestUserBySearch();
});

// Change user settings
Cypress.Commands.add('userSet', async data => {
  await login();
  await fetch(`${adminUrl}/users/${userId}`, {
    body: JSON.stringify({ ...user, ...data }),
    method: 'PUT',
    mode: 'cors',
    credentials: 'include',
  });
});

Cypress.Commands.add('getResetTokenURL', async () => {
  const resetTokenData = await fetch(`${adminUrl}/users/generate_new_reset_token`, {
    body: JSON.stringify({ user_id: userId, app_id: appId }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  }).then(r => r.json());
  return `${Cypress.config('baseUrl')}/password/reset/?appId=${appId}&url=${Cypress.config('serverUrl')}&token=${resetTokenData.Token}`;
});

Cypress.Commands.add('visitLogin', options => {
  return cy.visit(`${Cypress.config('baseUrl')}/login/?${new URLSearchParams({ ...options, appId: appId, url: Cypress.config('serverUrl') }).toString()}`);
});
Cypress.Commands.add('loginWithEmail', (email = 'test@test.com', password = 'Password') => {
  cy.get('[placeholder=Email]').click().type(email);
  cy.get('[placeholder=Password]').click().type(password);
  cy.screenshot();
  cy.get('button').contains('Login').click();
});

Cypress.Commands.add('registerWithEmail', (email = 'test@test.com', password = 'Password') => {
  cy.contains('Sign Up').click();
  cy.get('[placeholder=Email]').click().type(email);
  cy.get('[placeholder=Password]').click().type(password);
  cy.contains('Go back to login');
  cy.screenshot();
  cy.get('button').contains('Continue').click();
});

Cypress.Commands.add('verifyTfa', (code = '0000') => {
  cy.get('[placeholder="Verify code"]').click().type('0000');
  cy.contains('Go back to login');
  cy.screenshot();
  cy.contains('Confirm').click();
});
