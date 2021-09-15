import {
  applyRules, hasError, matches, notEmpty, longerThan, hasUpperLetter, emailFormat,
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
  username: [notEmpty('User name should not be empty')],
  fullName: [],
  password: [
    notEmpty('Password should not be empty'),
    longerThan(6, 'Password should have length of at least six characters'),
    hasUpperLetter('Password should contain at least one uppercase letter'),
  ],
  email: [
    notEmpty('Email should not be empty'),
    emailFormat('Email format is invalid'),
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
  const errors = validate('all', values);

  return hasError(errors) ? errors : {};
};

export default rules;
