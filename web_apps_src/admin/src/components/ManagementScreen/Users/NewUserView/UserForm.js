import { hasError } from '@dprovodnikov/validation';
import PropTypes from 'prop-types';
import React from 'react';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import Button from '~/components/shared/Button';
import Field from '~/components/shared/Field';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import MultipleInput from '~/components/shared/MultipleInput';
import Toggle from '~/components/shared/Toggle';
import useForm from '~/hooks/useForm';
import { toDeepCase } from '~/utils/apiMapper';
import './UserForm.css';
import { validateUserForm } from './validation';


const UserForm = ({ saving, error, onCancel, onSubmit }) => {
  const initialValues = {
    username: '',
    password: '',
    fullName: '',
    email: '',
    phone: '',
    scopes: [],
    confirmPassword: '',
    tfaEnabled: false,
    role: '',
  };

  const handleSubmit = (values) => {
    onSubmit(toDeepCase({
      username: values.username,
      pswd: values.password,
      fullName: values.fullName,
      scopes: values.scopes,
      accessRole: values.role,
      email: values.email,
      phone: values.phone,
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

      <Field label="Full name">
        <Input
          name="fullName"
          value={form.values.fullName}
          placeholder="Enter full name"
          onChange={form.handleChange}
          errorMessage={form.errors.fullName}
        />
      </Field>

      <Field label="Email">
        <Input
          name="email"
          value={form.values.email}
          placeholder="Enter user email"
          onChange={form.handleChange}
          errorMessage={form.errors.email}
        />
      </Field>

      <Field label="Phone">
        <Input
          name="phone"
          value={form.values.phone}
          placeholder="Enter user phone number"
          onChange={form.handleChange}
          errorMessage={form.errors.phone}
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

      <Field label="user scopes">
        <MultipleInput
          values={form.values.scopes}
          placeholder="Hit Enter to add scope"
          onChange={s => form.setValue('scopes', s)}
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
