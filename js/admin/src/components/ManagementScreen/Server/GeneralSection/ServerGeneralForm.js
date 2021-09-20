import React from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import useForm from '~/hooks/useForm';
import MultipleInput from '~/components/shared/MultipleInput';

const GeneralForm = (props) => {
  const { error, settings, loading, onSubmit } = props;

  const initialState = {
    host: '',
    issuer: '',
    port: '',
    supportedScopes: [],
  };

  const handleSubmit = (values) => {
    onSubmit(update(settings, values));
  };

  const form = useForm(initialState, null, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues({
      host: settings.host,
      issuer: settings.issuer,
      port: settings.port,
      supportedScopes: settings.supportedScopes || [],
    });
  }, [settings]);

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Host">
        <Input
          name="host"
          value={form.values.host}
          autoComplete="off"
          placeholder="Enter host url"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Issuer">
        <Input
          name="issuer"
          value={form.values.issuer}
          autoComplete="off"
          placeholder="Enter issuer url"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Port">
        <Input
          name="port"
          value={form.values.port}
          autoComplete="off"
          placeholder="Enter port"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="supported scopes">
        <MultipleInput
          values={form.values.supportedScopes}
          placeholder="Hit Enter to add scope"
          onChange={s => form.setValue('supportedScopes', s)}
        />
      </Field>

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          Icon={loading ? LoadingIcon : SaveIcon}
          disabled={loading}
          error={!loading && !!error}
        >
          Save Changes
        </Button>
      </footer>
    </form>
  );
};

export default GeneralForm;
