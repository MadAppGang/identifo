import {
  applyRules, hasError, notEmpty,
} from '@dprovodnikov/validation';


const schema = {
  folder: [
    notEmpty('Field is required'),
  ],
  bucket: [
    notEmpty('Field is required'),
  ],
  region: [
    notEmpty('Field is required'),
  ],
};

export const validateForm = (values) => {
  const validate = applyRules(schema);

  const omitFieldsByStorageType = {
    local: ['bucket', 'region'],
    none: ['folder', 'bucket', 'region'],
    default: ['folder', 'bucket', 'region'],
  };
  const errors = validate('all', values[values.type], {
    omit: omitFieldsByStorageType[values.type],
  });

  return hasError(errors) ? errors : {};
};
