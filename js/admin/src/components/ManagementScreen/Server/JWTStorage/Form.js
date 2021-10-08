import update from '@madappgang/update-by-path';
import React from 'react';
import CheckIcon from '~/components/icons/CheckIcon.svg';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import Button from '~/components/shared/Button';
import Field from '~/components/shared/Field';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import { Option, Select } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';
import { verificationStatuses } from '~/enums';

const storageTypes = {
  FILE: 'file',
  LOCAL: 'local',
  S3: 's3',
};

const ServerJWTForm = (props) => {
  const { error, loading, settings,
    onSubmit, verificationStatus, onChange, handleVerify } = props;
  const initialValues = {
    type: '',
    [storageTypes.FILE]: {
      privateKeyPath: undefined,
    },
    [storageTypes.S3]: {
      region: undefined,
      bucket: undefined,
      publicKeyKey: undefined,
      privateKeyKey: undefined,
    },
  };

  const handleSubmit = (values) => {
    onSubmit(update(settings, values));
  };
  // TODO: Nikita K add validations
  const form = useForm(initialValues, null, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(settings);
  }, [settings]);
  React.useEffect(() => {
    if (JSON.stringify(settings) !== JSON.stringify(form.values)) {
      onChange();
    }
  }, [settings, form.values]);

  if (!form) return null;

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Storage Type">
        <Select
          value={form.values.type}
          disabled={loading || !form.values.type}
          onChange={value => form.setValue('type', value)}
          placeholder="Select storage type"
        >
          <Option value={storageTypes.LOCAL} title="Local" />
          <Option value={storageTypes.S3} title="S3" />
        </Select>
      </Field>

      {form.values.type === storageTypes.S3 && (
        <>
          <Field label="Region">
            <Input
              name={`${storageTypes.S3}.region`}
              value={form.values[storageTypes.S3].region}
              autoComplete="off"
              placeholder="Enter s3 region"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Bucket" subtext="Can be overriden by IDENTIFO_JWT_KEYS_BUCKET env variable">
            <Input
              name={`${storageTypes.S3}.bucket`}
              value={form.values[storageTypes.S3].bucket}
              autoComplete="off"
              placeholder="Enter s3 bucket"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Public key S3 key">
            <Input
              name={`${storageTypes.S3}.publicKeyKey`}
              value={form.values[storageTypes.S3].publicKeyKey}
              autoComplete="off"
              placeholder="Enter s3 public key"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Private key S3 key">
            <Input
              name={`${storageTypes.S3}.privateKeyKey`}
              value={form.values[storageTypes.S3].privateKeyKey}
              autoComplete="off"
              placeholder="Enter s3 private key"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
        </>
      )}
      {form.values.type === storageTypes.LOCAL && (
        <>
          <Field label="Private key path">
            <Input
              name={`${storageTypes.FILE}.privateKeyPath`}
              value={form.values[storageTypes.FILE].privateKeyPath}
              autoComplete="off"
              placeholder="Enter private key path"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
        </>
      )}
      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          Icon={loading ? LoadingIcon : SaveIcon}
          disabled={loading}
          error={!loading && !!error}
        >
          Save Changes
        </Button>
        <Button
          error={verificationStatus === verificationStatuses.fail}
          success={verificationStatus === verificationStatuses.success}
          outline={verificationStatus === verificationStatuses.required}
          type="button"
          onClick={() => handleVerify(update(settings, form.values))}
          Icon={loading ? LoadingIcon : CheckIcon}
          disabled={loading}
        >
            Verify
        </Button>
      </footer>
    </form>
  );
};

export default ServerJWTForm;
