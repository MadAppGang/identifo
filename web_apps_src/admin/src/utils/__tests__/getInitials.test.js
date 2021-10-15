import { getInitials } from '..';

describe('getInitials util function', () => {
  test('returns correct initials when name is present', () => {
    expect(getInitials('john doe')).toBe('JD');
  });

  test('returns first two letters of email when name is absent', () => {
    expect(getInitials('', 'johndoe@gmail.com')).toBe('JO');
  });

  test('returns first two letters of first name is last name is absent', () => {
    expect(getInitials('john')).toBe('JO');
  });
});
