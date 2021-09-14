import {
  applyRules, hasError, notEmpty,
} from '@dprovodnikov/validation';

const onlyDigits = message => (value) => {
  if (!value) {
    return '';
  }

  if (!Number(value)) {
    return message;
  }

  return Number.isNaN(Number(value)) ? message : '';
};

export const adminAccountFormRules = {
  loginEnvName: [
    notEmpty('Login env name should not be empty'),
  ],
  passwordEnvName: [
    notEmpty('Password env name should not be empty'),
  ],
  sessionDuration: [
    notEmpty('You have to specify session duration (in seconds)'),
    onlyDigits('Duration should be specified in seconds'),
  ],
};

export const sessionStorageFormRules = {
  sessionDuration: [
    notEmpty('You have to specify session duration (in seconds)'),
    onlyDigits('Duration should be specified in seconds'),
  ],
  address: [
    notEmpty('You have to specify address'),
  ],
  password: [
    notEmpty('You have to specify password'),
  ],
  db: [
    notEmpty('You have to specify db of type int'),
    onlyDigits('Db should be of type int'),
  ],
  region: [
    notEmpty('You have to specify region'),
  ],
};


export const validateAccountForm = (values) => {
  const validate = applyRules(adminAccountFormRules);
  const errors = validate('all', values);
  return hasError(errors) ? errors : {};
};

export const validateSessionStorageForm = (values) => {
  const validate = applyRules(sessionStorageFormRules);

  const omitFieldsByStorageType = {
    memory: ['address', 'password', 'db', 'region', 'endpoint'],
    redis: ['region', 'endpoint'],
    dynamodb: ['address', 'password', 'db'],
  };
  const errors = validate('all', values, {
    omit: omitFieldsByStorageType[values.type],
  });

  return hasError(errors) ? errors : {};
};
