import React from 'react';
import { hasError } from '@dprovodnikov/validation';
import PropTypes from 'prop-types';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import Toggle from '~/components/shared/Toggle';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import useForm from '~/hooks/useForm';
import { validateUserForm } from './validation';
import { toDeepCase } from '~/utils/apiMapper';

import './UserForm.css';

const UserForm = ({ saving, error, onCancel, onSubmit }) => {
  const initialValues = {
    username: '',
    password: '',
    confirmPassword: '',
    tfaEnabled: false,
    role: '',
  };

  const handleSubmit = (values) => {
    onSubmit(toDeepCase({
      username: values.username,
      password: values.password,
      accessRole: values.role,
      tfaInfo: {
        isEnabled: values.tfaEnabled,
      },
    }, 'snake'));
  };

  const form = useForm(initialValues, validateUserForm, handleSubmit);

  return (
    <form className="iap-users-form" onSubmit={form.handleSubmit}>
      {error && <FormErrorMessage error={error} />}

      <Field label="Username">
        <Input
          name="username"
          value={form.values.username}
          placeholder="Enter username"
          onChange={form.handleChange}
          errorMessage={form.errors.username}
        />
      </Field>

      <Field label="Access Role">
        <Input
          name="role"
          value={form.values.role}
          placeholder="Enter access role"
          onChange={form.handleChange}
        />
      </Field>

      <Field label="Password">
        <Input
          name="password"
          type="password"
          value={form.values.password}
          placeholder="Enter password"
          onChange={form.handleChange}
          errorMessage={form.errors.password}
        />
      </Field>

      <Field label="Confirm Password">
        <Input
          name="confirmPassword"
          type="password"
          value={form.values.confirmPassword}
          placeholder="Enter password once more"
          onChange={form.handleChange}
          errorMessage={form.errors.confirmPassword}
        />
      </Field>

      <Toggle
        label="Enable 2FA"
        value={form.values.tfaEnabled}
        onChange={value => form.setValue('tfaEnabled', value)}
      />

      <footer className="iap-users-form__footer">
        <Button
          type="submit"
          Icon={saving ? LoadingIcon : SaveIcon}
          disabled={saving || hasError(form.errors)}
          error={!saving && !!error}
        >
          Save User
        </Button>
        <Button transparent onClick={onCancel} disabled={saving}>
          Cancel
        </Button>
      </footer>
    </form>
  );
};

UserForm.propTypes = {
  onCancel: PropTypes.func,
  onSubmit: PropTypes.func,
  saving: PropTypes.bool,
  error: PropTypes.instanceOf(Error),
};

UserForm.defaultProps = {
  onCancel: () => {},
  onSubmit: () => {},
  saving: false,
  error: null,
};

export default UserForm;
