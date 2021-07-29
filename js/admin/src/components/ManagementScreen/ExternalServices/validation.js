import { applyRules, hasError, notEmpty } from '@dprovodnikov/validation';

export const smsServiceValidationRules = {
  type: [notEmpty('You should select sms service')],
  accountSid: [notEmpty('Account SID cannot be empty')],
  authToken: [notEmpty('Token cannot be empty')],
  serviceSid: [notEmpty('Service SID cannot be empty')],
};

const validate = applyRules(smsServiceValidationRules);

export const validateSmsServiceForm = (values) => {
  const errors = validate('all', values, {
    omit: values.type === 'mock' ? ['accountSid', 'authToken', 'serviceSid'] : [],
  });

  return hasError(errors) ? errors : {};
};
