import { pause } from '~/utils';

/* eslint-disable */

const data = {
  users: [
    {
      id: '507f1f77bcf86cd799439011',
      username: 'dprovodnikov',
      email: '@madappgang.com',
      latest_login_time: Date.now() / 1000,
      num_of_logins: 0,
      tfa_info: {
        is_enabled: true,
      },
      access_role: '',
      active: false,
      phone: '+380993233430',
    },
    {
      id: '507f1f77bcf86cd799439012',
      tfa_info: {
        is_enabled: false,
      },
      active: true,
      username: 'test3@madappgang.com',
      email: 'test3@madappgang.com',
      phone: '+380993233430',
    },
    {
      id: '507f1f77bcf86cd799439013',
      tfa_info: {
        is_enabled: false,
      },
      active: true,
      username: '',
      phone: '+380993233430',
    },
  ],
};

const createUserServiceMock = () => {
  const fetchUsers = async (filters = {}) => {
    const { search } = filters;

    if (search) {
      await pause(300);

      return data.users.filter((user) => {
        return user.name.toLowerCase().includes(search.toLowerCase())
        || user.email.toLowerCase().includes(search.toLowerCase());
      });
    }

    await pause(600);

    return {
      users: data.users,
      total: data.users.length,
    };
  };

  const postUser = async (user) => {
    await pause(600);

    if (user.name === 'Trigger Error') {
      throw new Error('This name is already taken.');
    }

    const insertion = {
      ...user,
      id: Date.now().toString(),
      num_of_logins: 0,
      latest_login_time: undefined,
      password: undefined,
    };

    data.users.push(insertion);

    return insertion;
  };

  const alterUser = async (id, changes = {}) => {
    await pause(600);

    if (changes.name === 'Trigger Error') {
      throw new Error('This name is already taken');
    }

    data.users = data.users.map((user) => {
      if (user.id === id) {
        return {
          ...user,
          ...changes,
        };
      }
      return user;
    });

    return data.users.find(user => user.id === id);
  };

  const fetchUserById = async (id) => {
    const user = data.users.find(u => u.id === id);

    if (!user) {
      return Promise.reject(new Error('User is not found'));
    }

    return user;
  };

  const deleteUserById = async (id) => {
    await pause(600);

    data.users = data.users.filter(u => u.id !== id);
  };

  return Object.freeze({
    fetchUsers,
    postUser,
    alterUser,
    fetchUserById,
    deleteUserById,
  });
};

export default createUserServiceMock;
