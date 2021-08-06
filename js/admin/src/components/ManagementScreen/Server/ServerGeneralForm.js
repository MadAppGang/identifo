import React from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';
import { domChangeEvent } from '~/utils';

const GeneralForm = (props) => {
  const { error, settings, loading, onSubmit } = props;

  const initialState = {
    host: settings ? settings.host : '',
    issuer: settings ? settings.issuer : '',
    algorithm: settings ? settings.algorithm : '',
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
      algorithm: settings.algorithm,
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

      <Field label="Algorithm">
        <Select
          value={form.values.algorithm}
          disabled={loading}
          onChange={v => form.handleChange(domChangeEvent('algorithm', v))}
          placeholder="Select Algorithm"
        >
          <Option value="auto" title="Auto" />
          <Option value="rs256" title="rs256" />
          <Option value="es256" title="es256" />
        </Select>
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
