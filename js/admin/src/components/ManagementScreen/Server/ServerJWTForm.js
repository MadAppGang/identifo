import React from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import FileInput from '~/components/shared/FileInput';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';
import { domChangeEvent } from '~/utils';

const storageTypes = {
  FILE: 'file',
  S3: 's3',
};

const validateJwtForm = (values) => {
  const errors = {};

  if (values.privateKeyFile && !values.publicKeyFile) {
    errors.publicKeyFile = 'Both keys are required';
  }

  if (values.publicKeyFile && !values.privateKeyFile) {
    errors.privateKeyFile = 'Both keys are required';
  }

  return errors;
};

const ServerJWTForm = (props) => {
  const { error, loading, onSubmit } = props;

  const settings = props.settings ? props.settings.keyStorage : null;

  const initialValues = {
    storageType: settings ? settings.type : '',
    publicKeyPath: settings ? settings.publicKey : '',
    privateKeyPath: settings ? settings.privateKey : '',
    region: settings ? settings.region : '',
    bucket: settings ? settings.bucket : '',
    publicKeyFile: null,
    privateKeyFile: null,
  };

  const handleSubmit = (values) => {
    onSubmit(update(props.settings, {
      keyStorage: {
        type: values.storageType,
        publicKey: values.publicKeyPath,
        privateKey: values.privateKeyPath,
        region: values.region,
        bucket: values.bucket,
      },
      publicKey: values.publicKeyFile,
      privateKey: values.privateKeyFile,
    }));
  };

  const form = useForm(initialValues, validateJwtForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(update(form.values, {
      storageType: settings.type,
      publicKeyPath: settings.publicKey,
      privateKeyPath: settings.privateKey,
      region: settings.region,
      bucket: settings.bucket,
    }));
  }, [settings]);

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Storage Type">
        <Select
          value={form.values.storageType}
          disabled={loading}
          onChange={value => form.setValue('storageType', value)}
          placeholder="Select storage type"
        >
          <Option value={storageTypes.FILE} title="File" />
          <Option value={storageTypes.S3} title="S3" />
        </Select>
      </Field>

      {form.values.storageType === storageTypes.S3 && (
        <Field label="Region">
          <Input
            name="region"
            value={form.values.region}
            autoComplete="off"
            placeholder="Enter s3 region"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      {form.values.storageType === storageTypes.S3 && (
        <Field label="Bucket" subtext="Can be overriden by IDENTIFO_JWT_KEYS_BUCKET env variable">
          <Input
            name="bucket"
            value={form.values.bucket}
            autoComplete="off"
            placeholder="Enter s3 bucket"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      <Field
        label="Public Key"
        subtext={form.values.publicKeyFile ? form.values.publicKeyFile.name : 'No file selected'}
      >
        <FileInput
          path={form.values.publicKeyPath}
          placeholder="Specify path to folder"
          onFile={file => form.setValue('publicKeyFile', file)}
          onPath={v => form.handleChange(domChangeEvent('publicKeyPath', v))}
          errorMessage={form.errors.publicKeyFile}
          disabled={loading}
        />
      </Field>

      <Field
        label="Private Key"
        subtext={form.values.privateKeyFile ? form.values.privateKeyFile.name : 'No file selected'}
      >
        <FileInput
          path={form.values.privateKeyPath}
          placeholder="Specify path to folder"
          onFile={file => form.setValue('privateKeyFile', file)}
          onPath={v => form.handleChange(domChangeEvent('privateKeyPath', v))}
          errorMessage={form.errors.privateKeyFile}
          disabled={loading}
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

export default ServerJWTForm;
