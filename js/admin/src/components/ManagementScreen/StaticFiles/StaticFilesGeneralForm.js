import React from 'react';
import update from '@madappgang/update-by-path';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';
import Toggle from '~/components/shared/Toggle';

const [LOCAL, S3, DYNAMO_DB] = ['local', 's3', 'dynamo'];

const StaticFilesGeneralForm = (props) => {
  const { loading, error, settings, onSubmit } = props;

  const initialValues = {
    [LOCAL]: { folder: '' },
    [S3]: { region: '', bucket: '', folder: '' },
    [DYNAMO_DB]: { region: '', endpoint: '' },
    serveAdminPanel: false,
    type: settings.type || LOCAL,
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
      <Field label="Storage Type">
        <Select
          value={form.values.type || ''}
          onChange={value => form.setValue('type', value)}
          placeholder="Select storage type"
        >
          <Option value={LOCAL} title="Local" />
          <Option value={S3} title="S3" />
          <Option value={DYNAMO_DB} title="DynamoDB" />
        </Select>
      </Field>

      {form.values.type === LOCAL && (
        <Field label="Folder">
          <Input
            name="local.folder"
            value={form.values.local.folder}
            autoComplete="off"
            placeholder="Specify folder"
            onChange={form.handleChange}
            disabled={loading}
          />
        </Field>
      )}

      {form.values.type === S3 && (
        <>
          <Field label="Region">
            <Input
              name="s3.region"
              value={form.values.s3.region}
              autoComplete="off"
              placeholder="Specify region"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Bucket">
            <Input
              name="s3.bucket"
              value={form.values.s3.bucket}
              autoComplete="off"
              placeholder="Specify bucket"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Folder">
            <Input
              name="s3.folder"
              value={form.values.s3.folder}
              autoComplete="off"
              placeholder="Specify folder"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
        </>
      )}

      {form.values.type === DYNAMO_DB && (
        <>
          <Field label="Region">
            <Input
              name="dynamo.region"
              value={form.values.dynamo.region}
              autoComplete="off"
              placeholder="Specify db region"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
          <Field label="Endpoint" subtext="Can be omitted if region is set">
            <Input
              name="dynamo.endpoint"
              value={form.values.dynamo.endpoint}
              autoComplete="off"
              placeholder="Specify db endpoint"
              onChange={form.handleChange}
              disabled={loading}
            />
          </Field>
        </>
      )}
      <Field label="Serve admin panel">
        <div className="iap-apps-form--toggler">
          <Toggle value={form.values.serveAdminPanel} onChange={v => form.setValue('serveAdminPanel', v)} />
        </div>
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

export default StaticFilesGeneralForm;
