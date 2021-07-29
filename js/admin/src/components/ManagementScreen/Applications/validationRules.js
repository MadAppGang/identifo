import { notEmpty, urlFormat } from '@dprovodnikov/validation';

const redirectUrlRule = message => (value) => {
  if (!value) {
    return '';
  }

  if (Array.isArray(value)) {
    const hasError = value.some(v => !!urlFormat(message)(v));

    if (hasError) {
      return message;
    }

    return '';
  }

  return urlFormat(message)(value);
};

const onlyDigits = message => (value) => {
  if (!value) {
    return '';
  }

  if (!Number(value)) {
    return message;
  }

  return Number.isNaN(Number(value)) ? message : '';
};

const rules = {
  name: [notEmpty('Application name should not be empty')],
  type: [notEmpty('Application type should be selected')],
  redirectUrls: [
    redirectUrlRule('Url format is invalid. Make sure you include scheme.'),
  ],
  tokenLifespan: [
    onlyDigits('Token lifespan can only contain digits'),
  ],
};

export default rules;
