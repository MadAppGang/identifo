import user from '../fixtures/user.json';
import app from '../fixtures/app.json';
import server from '../fixtures/server.json';
const adminUrl = `${Cypress.config('serverUrl')}/admin`;
let appId = [''];
let lastAppId = '';
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
  if (!userId) {
    return;
  }
  return fetch(`${adminUrl}/users/${userId}`, {
    body: null,
    method: 'DELETE',
    mode: 'cors',
    credentials: 'include',
  });
};
const deleteTestApp = async () => {
  if (!appId) {
    return;
  }
  const p = [];
  for (const t of appId) {
    p.push(
      fetch(`${adminUrl}/apps/${t}`, {
        body: null,
        method: 'DELETE',
        mode: 'cors',
        credentials: 'include',
      }),
    );
  }

  return Promise.all(p);
};
const parseJwt = token => {
  var base64Url = token.split('.')[1];
  var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  var jsonPayload = decodeURIComponent(
    atob(base64)
      .split('')
      .map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join(''),
  );

  return JSON.parse(jsonPayload);
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
    .then(result => {
      appId.push(result.id);
      lastAppId = result.id;
    });
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
Cypress.Commands.add('createApp', async oApp => {
  await login();
  await fetch(`${adminUrl}/apps`, {
    body: JSON.stringify({
      ...app,
      ...oApp,
    }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => {
      appId.push(result.id);
      lastAppId = result.id;
    });
});
Cypress.Commands.add('createUser', async u => {
  await login();
  await fetch(`${adminUrl}/users/`, {
    body: JSON.stringify({ ...user, ...u }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  })
    .then(r => r.json())
    .then(result => (userId = result.id));
});
Cypress.Commands.add('deleteUser', async searchString => {
  const testUser = await fetch(`${adminUrl}/users?search=${searchString}`, {
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
});
// Cleanup test user and app
Cypress.Commands.add('deleteAppAndUser', async data => {
  await login();
  await deleteTestApp();
  await deleteTestUser();
  await deleteTestUserBySearch();
});
// Change app settings
Cypress.Commands.add('appSet', async data => {
  await login();
  await fetch(`${adminUrl}/apps/${lastAppId}`, {
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
  await new Promise(resolve => setTimeout(resolve, 2000));
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
  await login();
  console.log(userId);
  const resetTokenData = await fetch(`${adminUrl}/users/generate_new_reset_token`, {
    body: JSON.stringify({ user_id: userId, app_id: lastAppId }),
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  }).then(r => r.json());
  if (resetTokenData.URL.indexOf('password/reset') === -1) {
    throw new Error('Invalid reset token URL');
  }
  return `${Cypress.config('baseUrl')}/password/reset/?appId=${lastAppId}&url=${Cypress.config('serverUrl')}&token=${resetTokenData.Token}`;
});

Cypress.Commands.add('addInvite', async (email, role = 'user') => {
  await login();
  const invite = await fetch(`${adminUrl}/invites`, {
    body: `{"email":"${email}", "app_id":"${lastAppId}", "access_role":"${role}"}`,
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
  }).then(r => r.json());
  return `${Cypress.config('baseUrl')}/register/?email=${email}&appId=${lastAppId}&token=${invite.token}`;
});

Cypress.Commands.add('visitLogin', options => {
  window.localStorage.setItem('debug', true);
  return cy.visit(`${Cypress.config('baseUrl')}/login/?${new URLSearchParams({ ...options, appId: lastAppId, url: Cypress.config('serverUrl') }).toString()}`);
});
Cypress.Commands.add('loginWithEmail', p => {
  const login = { ...{ email: 'test@test.com', password: 'Password', remember: false }, ...p };
  cy.get('[placeholder=Email]').click().type(login.email);
  cy.get('[placeholder=Password]').click().type(login.password);
  if (login.remember) {
    cy.contains('Remember me').click();
  }
  cy.screenshot();
  cy.get('button').contains('Login').click();
});

Cypress.Commands.add('loginWithPhone', (phone = '+0123456789', remember = false) => {
  cy.get('[placeholder="Phone number"]').click().type(phone);
  if (remember) {
    cy.contains('Remember me').click();
  }
  cy.screenshot();
  cy.get('button').contains('Continue').click();
});

Cypress.Commands.add('loginWithPhoneVerify', (code = '0000') => {
  cy.get('[placeholder="Verify code"]').click().type('0000');
  cy.contains('Go back to login');
  cy.screenshot();
  cy.contains('Confirm').click();
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

Cypress.Commands.add('verifySuccessToken', () => {
  cy.contains('Success');
  cy.get('#access_token')
    .invoke('text')
    .then(text => {
      const token = parseJwt(text.trim());
      expect(token.scopes).equal(undefined);
      expect(token.type).equal('access');
    });
});

Cypress.Commands.add('verifyRefreshSuccessToken', () => {
  cy.contains('Success');
  cy.get('#access_token')
    .invoke('text')
    .then(text => {
      const token = parseJwt(text.trim());
      expect(token.scopes).equal('offline');
      expect(token.type).equal('access');
    });
  cy.get('#refresh_token')
    .invoke('text')
    .then(text => {
      const token = parseJwt(text.trim());
      expect(token.scopes).equal('offline');
      expect(token.type).equal('refresh');
    });
});
Cypress.Commands.add('getUserData', () => {
  return cy
    .get('#user_data')
    .invoke('text')
    .then(text => {
      return JSON.parse(text.trim());
    });
});
