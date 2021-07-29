import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import MultipleInput from '~/components/shared/MultipleInput';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';

const storageTypes = {
  FILE: 'file',
  S3: 's3',
  ETCD: 'etcd',
};

const settingsKeyDescription = {
  [storageTypes.FILE]: 'Must be a filename (e.g, server-config.yaml)',
  [storageTypes.ETCD]: 'Must be a key name (e.g, identifo/server-settings)',
  [storageTypes.S3]: 'Must be a name of an object in the bucket (e.g, server-config.yaml)',
};

const ServerConfigurationForm = (props) => {
  const { loading, error, settings, onSubmit } = props;

  const [storageType, setStorageType] = useState(settings ? settings.type : '');
  const [settingsKey, setSettingsKey] = useState(settings ? settings.settingsKey : '');
  const [endpoints, setEndpoints] = useState(settings ? settings.endpoints : []);
  const [region, setRegion] = useState(settings ? settings.region : '');
  const [bucket, setBucket] = useState(settings ? settings.bucket : '');

  useEffect(() => {
    if (!settings) return;

    setStorageType(settings.type || '');
    setSettingsKey(settings.settingsKey || '');
    setEndpoints(settings.endpoints || []);
    setRegion(settings.region || '');
    setBucket(settings.bucket || '');
  }, [settings]);

  const handleSubmit = (event) => {
    event.preventDefault();
    onSubmit(update(settings, {
      type: storageType, settingsKey, endpoints, region, bucket,
    }));
  };

  return (
    <form className="iap-apps-form" onSubmit={handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Storage Type">
        <Select
          value={storageType}
          disabled={loading}
          onChange={setStorageType}
          placeholder="Select storage type"
        >
          <Option value={storageTypes.FILE} title="File" />
          <Option value={storageTypes.ETCD} title="Etcd" />
          <Option value={storageTypes.S3} title="S3" />
        </Select>
      </Field>

      <Field label="Settings Key" subtext={settingsKeyDescription[storageType]}>
        <Input
          value={settingsKey}
          autoComplete="off"
          placeholder="Enter settings key"
          onChange={e => setSettingsKey(e.target.value)}
          disabled={loading}
        />
      </Field>

      {storageType === storageTypes.ETCD && (
        <Field label="Etcd Endpoints">
          <MultipleInput
            values={endpoints}
            placeholder="Hit Enter to add endpoint"
            onChange={setEndpoints}
          />
        </Field>
      )}

      {storageType === storageTypes.S3 && (
        <Field label="Region">
          <Input
            value={region}
            autoComplete="off"
            placeholder="Enter s3 region"
            onValue={setRegion}
            disabled={loading}
          />
        </Field>
      )}

      {storageType === storageTypes.S3 && (
        <Field
          label="Bucket"
          subtext="Can be overriden by IDENTIFO_JWT_KEYS_BUCKET env variable"
        >
          <Input
            value={bucket}
            autoComplete="off"
            placeholder="Enter s3 bucket"
            onValue={setBucket}
            disabled={loading}
          />
        </Field>
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
      </footer>
    </form>
  );
}

ServerConfigurationForm.defaultProps = {
  settings: {},
  loading: false,
  error: null,
  onSubmit: () => null,
};

export default ServerConfigurationForm;
