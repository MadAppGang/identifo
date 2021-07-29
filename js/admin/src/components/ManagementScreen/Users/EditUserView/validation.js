import {
  hasError, applyRules, matches, notEmpty, emailFormat, longerThan, hasUpperLetter,
} from '@dprovodnikov/validation';

const phoneNumberRegExp = /^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$/;

const phoneNumberRule = message => (value) => {
  if (!value) {
    return '';
  }

  if (!phoneNumberRegExp.test(value)) {
    return message;
  }

  return '';
};

const rules = {
  username: [notEmpty('Username should not be empty')],
  email: [
    notEmpty('Email should not be empty'),
    emailFormat('Email format is invalid'),
  ],
  password: [
    notEmpty('Password should not be empty'),
    longerThan(6, 'Password should have length of at least six characters'),
    hasUpperLetter('Password should contain at least one uppercase letter'),
  ],
  confirmPassword: [
    notEmpty('You should confirm password'),
    matches('password', 'Passwords do not match'),
  ],
  phone: [
    phoneNumberRule('Phone number format is invalid'),
  ],
};

const validate = applyRules(rules);

export const validateUserForm = (values) => {
  const errors = validate('all', values, {
    omit: values.editPassword ? [] : ['password', 'confirmPassword'],
  });

  return hasError(errors) ? errors : {};
};
