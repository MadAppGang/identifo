import {
  applyRules, hasError, matches, notEmpty, longerThan, hasUpperLetter,
} from '@dprovodnikov/validation';

const rules = {
  username: [notEmpty('Username should not be empty')],
  password: [
    notEmpty('Password should not be empty'),
    longerThan(6, 'Password should have length of at least six characters'),
    hasUpperLetter('Password should contain at least one uppercase letter'),
  ],
  confirmPassword: [
    notEmpty('You should confirm password'),
    matches('password', 'Passwords do not match'),
  ],
};

const validate = applyRules(rules);

export const validateUserForm = (values) => {
  const errors = validate('all', values);

  return hasError(errors) ? errors : {};
};

export default rules;
