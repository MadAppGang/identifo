import { notEmpty } from '@dprovodnikov/validation';

const rules = {
  type: [notEmpty('Type should not be empty')],
  region: [notEmpty('Region should not be empty')],
  name: [notEmpty('Name should not be empty')],
  endpoint: [notEmpty('Endpoint should not be empty')],
  path: [notEmpty('Path should not be empty')],
};

export default rules;
