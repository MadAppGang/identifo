import React from 'react';
import update from '@madappgang/update-by-path';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import Input from '~/components/shared/Input';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';

const [MOCK, AWS_SES, MAILGUN] = ['mock', 'aws ses', 'mailgun'];

const MailServiceSettings = (props) => {
  const { loading, settings, onSubmit } = props;

  const initialValues = {
    type: settings ? settings.type : '',
    domain: settings ? settings.domain : '',
    privateKey: settings ? settings.privateKey : '',
    publicKey: settings ? settings.publicKey : '',
    sender: settings ? settings.sender : '',
    region: settings ? settings.region : '',
  };

  const handleSubmit = (values) => {
    onSubmit(update(settings, values));
  };

  const form = useForm(initialValues, null, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(settings);
  }, [settings]);

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      <Field label="Mail Service">
        <Select
          value={form.values.type}
          disabled={loading}
          onChange={value => form.setValue('type', value)}
          placeholder="Select Supported Service"
        >
          <Option value={MAILGUN} title="Mailgun" />
          <Option value={AWS_SES} title="Amazon SES" />
          <Option value={MOCK} title="Mock" />
        </Select>
      </Field>

      <Field
        label="Sender"
        subtext={'If can be overriden by "MAILGUN_SENDER" or "AWS_SES_SENDER" env vars.'}
      >
        <Input
          name="sender"
          value={form.values.sender}
          autoComplete="off"
          placeholder="Specify Sender"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      {form.values.type === MAILGUN && (
        <Field label="Domain" subtext="Can be overriden by MAILGUN_DOMAIN env var">
          <Input
            name="domain"
            value={form.values.domain}
            autoComplete="off"
            placeholder="Specify Mailgun domain"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      {form.values.type === MAILGUN && (
        <Field label="Public Key" subtext="Can be overriden by MAILGUN_PUBLIC_KEY env var">
          <Input
            name="publicKey"
            value={form.values.publicKey}
            autoComplete="off"
            placeholder="Specify Mailgun public key"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      {form.values.type === MAILGUN && (
        <Field label="Private Key" subtext="Can be overriden by MAILGUN_PRIVATE_KEY env var">
          <Input
            name="privateKey"
            value={form.values.privateKey}
            autoComplete="off"
            placeholder="Specify Mailgun private key"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      {form.values.type === AWS_SES && (
        <Field label="Region" subtext="Can be overriden by AWS_SES_REGION env var">
          <Input
            name="region"
            value={form.values.region}
            autoComplete="off"
            placeholder="Specify Region"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          disabled={loading}
          Icon={loading ? LoadingIcon : SaveIcon}
        >
          Save Changes
        </Button>
      </footer>
    </form>
  );
};

export default MailServiceSettings;
