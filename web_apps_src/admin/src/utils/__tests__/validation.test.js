import { reset, hasError, applyRules } from '../validation';

describe('reset validation util function', () => {
  test('resets input validation properly', () => {
    const validation = {
      field: 'Value',
    };

    const expectedValidation = {
      field: '',
    };

    expect(reset(validation)).toEqual(expectedValidation);
  });
});

describe('hasError validation util function', () => {
  test('detects error in validation object properly', () => {
    const validation = {
      validField: '',
      invalidField: 'error',
    };

    expect(hasError(validation)).toBe(true);
  });
});

describe('applyRules  validation util function', () => {
  test('is a function', () => {
    expect(applyRules).toBeInstanceOf(Function);
  });

  test('returns a function', () => {
    expect(applyRules({})).toBeInstanceOf(Function);
  });

  test('returned function uses applied rule properly', () => {
    const validate = () => 'Validation message';

    const rules = {
      field: [validate],
    };

    expect(applyRules(rules)('field', { field: 'value' })).toBe('Validation message');
  });
});
