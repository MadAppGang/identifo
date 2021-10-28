import update from '@madappgang/update-by-path';
import React from 'react';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import Button from '~/components/shared/Button';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import { Option, Select } from '~/components/shared/Select';
import { verificationStatuses, storageTypes } from '~/enums';
import useForm from '~/hooks/useForm';
import { validateForm } from './validationRules';
import CheckIcon from '~/components/icons/CheckIcon.svg';


const [NONE, DEFAULT, LOCAL, S3] = ['none', 'default', storageTypes.local, storageTypes.s3];

const touched = (initial, values) => {
  return JSON.stringify(initial) !== JSON.stringify(values);
};


export const StaticFilesForm = (props) => {
  const { loading, error, settings, verificationStatus, onVerify, onChange, onSubmit } = props;

  const initialValues = {
    [LOCAL]: { folder: '' },
    [S3]: { region: '', bucket: '', folder: '' },
    type: settings.type || NONE,
  };
  const handleSubmit = async (values) => {
    onSubmit(update(settings, values));
  };

  const form = useForm(initialValues, validateForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(settings);
  }, [settings]);

  React.useEffect(() => {
    onChange(form.values);
  }, [form.values]);

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      <Field label="Storage Type">
        <Select
          value={form.values.type || ''}
          onChange={value => form.setValue('type', value)}
          placeholder="Select storage type"
        >
          <Option value={NONE} title="None" />
          <Option value={DEFAULT} title="Default" />
          <Option value={LOCAL} title="Local" />
          <Option value={S3} title="S3" />
        </Select>
      </Field>

      {form.values.type === LOCAL && (
        <Field label="Folder">
          <Input
            name="local.folder"
            errorMessage={form.errors.folder}
            value={form.values.local.folder}
            autoComplete="off"
            placeholder="Specify path to the folder"
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
              errorMessage={form.errors.region}
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
              errorMessage={form.errors.bucket}
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
              errorMessage={form.errors.folder}
              autoComplete="off"
              placeholder="Specify folder"
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
          error={!loading && !!error}
          disabled={!touched(settings, form.values) || loading}
        >
          Save Changes
        </Button>
        {(form.values.type === LOCAL || form.values.type === S3) && (
        <Button
          error={verificationStatus === verificationStatuses.fail}
          success={verificationStatus === verificationStatuses.success}
          outline={verificationStatus === verificationStatuses.required}
          type="button"
          onClick={() => onVerify(form.values)}
          disabled={!touched(settings, form.values)}
          Icon={loading ? LoadingIcon : CheckIcon}
        >
             Verify
        </Button>
        )}
      </footer>
    </form>
  );
};
