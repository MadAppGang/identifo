import React from 'react';
import update from '@madappgang/update-by-path';
import { hasError } from '@dprovodnikov/validation';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import useForm from '~/hooks/useForm';
import { validateAccountForm } from './validationRules';

const DEFAULT_SESSION_DURATION = 300;

const AdminAccountForm = ({ onSubmit, error, loading, settings }) => {
  const initialValues = {
    loginEnvName: settings ? settings.loginEnvName : '',
    passwordEnvName: settings ? settings.passwordEnvName : '',
    sessionDuration: '',
  };

  const handleSubmit = (data) => {
    const payload = {
      ...data,
      sessionDuration: Number(data.sessionDuration),
    };
    onSubmit(update(settings, payload));
  };

  const form = useForm(initialValues, validateAccountForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues({
      loginEnvName: settings.loginEnvName || '',
      passwordEnvName: settings.passwordEnvName || '',
      sessionDuration: settings.sessionDuration.toString() || DEFAULT_SESSION_DURATION.toString(),
    });
  }, [settings]);

  return (
    <form className="iap-settings-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}
      <Field label="session duration">
        <Input
          name="sessionDuration"
          value={form.values.sessionDuration}
          placeholder="e.g 300"
          onChange={form.handleChange}
          onBlur={form.handleBlur}
          errorMessage={form.errors.sessionDuration}
          disabled={loading}
        />
      </Field>
      <Field label="Login env name">
        <Input
          name="loginEnvName"
          value={form.values.loginEnvName}
          placeholder="e.g IDENTIFO_ADMIN_LOGIN"
          onChange={form.handleChange}
          onBlur={form.handleBlur}
          errorMessage={form.errors.loginEnvName}
          disabled={loading}
        />
      </Field>
      <Field label="Password env name">
        <Input
          name="passwordEnvName"
          value={form.values.passwordEnvName}
          placeholder="e.g IDENTIFO_ADMIN_PASSWORD"
          onChange={form.handleChange}
          onBlur={form.handleBlur}
          errorMessage={form.errors.passwordEnvName}
          disabled={loading}
        />
      </Field>
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
