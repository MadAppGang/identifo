import { notEmpty } from '@dprovodnikov/validation';

const rules = {
  region: [notEmpty('Region should not be empty')],
  name: [notEmpty('Name should not be empty')],
  endpoint: [notEmpty('Endpoint should not be empty')],
  path: [notEmpty('Path should not be empty')],
  database: [notEmpty('Database name should not be empty')],
  connection: [notEmpty('Connection should not be empty')],
};

export default rules;
