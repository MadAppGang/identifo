import React from 'react';
import { hasError } from '@dprovodnikov/validation';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import Toggle from '~/components/shared/Toggle';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { validateUserForm } from './validation';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import useForm from '~/hooks/useForm';
import { toDeepCase } from '~/utils/apiMapper';

const EditUserForm = ({ user, error, loading, onSubmit, onCancel }) => {
  const initialValues = {
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
    tfaEnabled: false,
    role: '',
    phone: '',
    active: false,
    editPassword: false,
  };

  const handleSubmit = (values) => {
    onSubmit(toDeepCase({
      email: values.email,
      username: values.username,
      password: values.editPassword ? values.password : '',
      tfaInfo: {
        isEnabled: values.tfaEnabled,
      },
      accessRole: values.role,
      phone: values.phone,
      active: values.active,
    }, 'snake'));
  };

  const form = useForm(initialValues, validateUserForm, handleSubmit);

  React.useEffect(() => {
    if (!user) return;

    form.setValues({
      email: user.email,
      username: user.username,
      tfaEnabled: user.tfa_info ? user.tfa_info.is_enabled : false,
      role: user.access_role || '',
      active: user.active || false,
      phone: user.phone || '',
      editPassword: form.values.editPassword,
      password: form.values.password,
      confirmPassword: form.values.confirmPassword,
    });
  }, [user]);

  return (
    <form className="iap-users-form" onSubmit={form.handleSubmit}>
      {!!error && <FormErrorMessage error={error} />}

      <Field label="Username">
        <Input
          name="username"
          value={form.values.username}
          placeholder="Enter username"
          onChange={form.handleChange}
          errorMessage={form.errors.username}
          disabled={loading}
        />
      </Field>

      <Field label="Access Role">
        <Input
          name="role"
          value={form.values.role}
          placeholder="Enter access role"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Email">
        <Input
          name="email"
          value={form.values.email}
          placeholder="Enter email"
          onChange={form.handleChange}
          errorMessage={form.errors.email}
          disabled={loading}
        />
      </Field>

      <Field label="Pnone Number">
        <Input
          name="phone"
          value={form.values.phone}
          placeholder="Enter phone number"
          onChange={form.handleChange}
          errorMessage={form.errors.phone}
          disabled={loading}
        />
      </Field>

      <div>
        <Toggle
          label="Enable 2FA"
          value={form.values.tfaEnabled}
          onChange={value => form.setValue('tfaEnabled', value)}
        />

        <Toggle
          label="Active"
          value={form.values.active}
          onChange={value => form.setValue('active', value)}
        />

        <Toggle
          label="Edit Password"
          value={form.values.editPassword}
          onChange={value => form.setValue('editPassword', value)}
        />
      </div>

      {form.values.editPassword && (
        <>
          <Field label="Password">
            <Input
              name="password"
              type="password"
              placeholder="Enter new password"
              value={form.values.password}
              onChange={form.handleChange}
              errorMessage={form.errors.password}
              disabled={loading}
            />
          </Field>

          <Field label="Confirm Password">
            <Input
              name="confirmPassword"
              type="password"
              placeholder="Enter new password"
              value={form.values.confirmPassword}
              onChange={form.handleChange}
              errorMessage={form.errors.confirmPassword}
              disabled={loading}
            />
          </Field>
        </>
      )}

      <footer className="iap-users-form__footer">
        <Button
          type="submit"
          Icon={loading ? LoadingIcon : SaveIcon}
          disabled={loading || hasError(form.errors)}
          error={!loading && !!error}
        >
          Save Changes
        </Button>
        <Button transparent onClick={onCancel} disabled={loading}>
          Cancel
        </Button>
      </footer>
    </form>
  );
};

EditUserForm.defaultProps = {
  onCancel: () => {},
};

export default EditUserForm;
