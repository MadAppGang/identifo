/* checks if at least one field of validation object has error message */
const hasError = validation => Object.values(validation).some(value => !!value);

/* applies rules to given value */
const validate = (rules = [], value, values) => {
  const results = rules.map(applyRule => applyRule(value, values));
  const errorResults = results.filter(result => !!result);
  return errorResults.length ? errorResults[0] : '';
};

/* generates a function that is used to validate fields */
const applyRules = validationConfig => (field, value, config = {}) => {
  if (field === 'all') {
    const { omit = [] } = config;

    return Object.keys(validationConfig)
      .filter(key => !omit.includes(key))
      .map(f => ({ [f]: validate(validationConfig[f], value[f], value) }))
      .reduce((output, entry) => ({ ...output, ...entry }), {});
  }

  return validate(validationConfig[field], value[field], value);
};

/* resets validation objects by removing all error messages */
const reset = validation => Object.entries(validation)
  .map(entry => [entry[0], ''])
  .reduce((output, entry) => ({ ...output, [entry[0]]: entry[1] }), {});

const emailFormatRule = message => (email) => {
  const emailRegExp = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

  if (!emailRegExp.test(email)) {
    return message;
  }

  return '';
};

const notEmptyRule = message => (value = '') => {
  if (!value.trim()) {
    return message;
  }

  return '';
};

const matchesRule = (comparisonField, message) => (value, fields) => {
  if (value !== fields[comparisonField]) {
    return message;
  }

  return '';
};

export {
  hasError,
  applyRules,
  emailFormatRule,
  notEmptyRule,
  matchesRule,
  reset,
};
