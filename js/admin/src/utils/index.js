import randomstring from 'randomstring';

export const pause = timeout => new Promise(resolve => setTimeout(resolve, timeout));

export const domChangeEvent = (name, value) => ({ target: { name, value } });

export const isPhone = (value) => {
  return !!(/^[+][0-9]{9,15}$/.test(value));
};

export const getError = (axiosErr) => {
  if (axiosErr.response && axiosErr.response.data) {
    return new Error(axiosErr.response.data.error);
  }

  return axiosErr;
};

export const getStatus = (axiosErr) => {
  const defaultStatus = 404;

  if (axiosErr.response) {
    return axiosErr.response.status || defaultStatus;
  }

  return defaultStatus;
};

export const getInitials = (fullName, email) => {
  const firstNOf = n => str => str.slice(0, n);
  const firstTwoOf = firstNOf(2);
  const firstOneOf = firstNOf(1);

  const [firstName, lastName] = fullName.split(/\s/);

  if (!firstName && !lastName) {
    return firstTwoOf(email).toUpperCase();
  }

  if (!lastName) {
    return firstTwoOf(firstName).toUpperCase();
  }

  return `${firstOneOf(firstName)}${firstOneOf(lastName)}`.toUpperCase();
};

export const generateSecret = (length = 24) => {
  return randomstring.generate(length, { charset: 'hex' });
};
