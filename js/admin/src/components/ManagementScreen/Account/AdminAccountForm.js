import React from 'react';
import update from '@madappgang/update-by-path';
import { hasError } from '@dprovodnikov/validation';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import Toggle from '~/components/shared/Toggle';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import useForm from '~/hooks/useForm';
import { validateAccountForm } from './validationRules';

const AdminAccountForm = ({ onSubmit, error, loading, settings }) => {
  const initialValues = {
    email: settings ? settings.email : '',
    password: '',
    confirmPassword: '',
    editPassword: false,
  };

  const handleSubmit = ({ email, password, editPassword }) => {
    onSubmit(update(settings, {
      email, password: editPassword ? password : undefined,
    }));
  };

  const form = useForm(initialValues, validateAccountForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues({
      email: settings.email || '',
      password: '',
      confirmPassword: '',
      editPassword: form.values.editPassword || false,
    });
  }, [settings]);

  return (
    <form className="iap-settings-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Email">
        <Input
          name="email"
          value={form.values.email}
          placeholder="Enter your email"
          onChange={form.handleChange}
          onBlur={form.handleBlur}
          errorMessage={form.errors.email}
          disabled={loading}
        />
      </Field>

      <Toggle
        label="Edit password"
        value={form.values.editPassword}
        onChange={(value) => {
          form.setErrors({
            password: '',
            confirmPassword: '',
          });

          form.setValues(update(form.values, {
            editPassword: value,
            password: '',
            confirmPassword: '',
          }));
        }}
      />

      {form.values.editPassword && (
        <div className="iap-settings-form__password-fields">
          <Field label="Password">
            <Input
              type="password"
              name="password"
              value={form.values.password}
              placeholder="Enter your password"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              errorMessage={form.errors.password}
            />
          </Field>

          <Field label="Confirm Password">
            <Input
              type="password"
              name="confirmPassword"
              value={form.values.confirmPassword}
              placeholder="Confirm your password"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              errorMessage={form.errors.confirmPassword}
            />
          </Field>
        </div>
      )}

      <footer className="iap-settings-form__footer">
        <Button
          type="submit"
          error={!loading && !!error}
          disabled={loading || hasError(form.errors)}
          Icon={loading ? LoadingIcon : SaveIcon}
        >
          Save Changes
        </Button>
      </footer>
    </form>
  );
};

export default AdminAccountForm;
