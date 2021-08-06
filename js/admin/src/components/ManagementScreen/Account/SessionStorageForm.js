import React from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';
import { validateSessionStorageForm } from './validationRules';
import useForm from '~/hooks/useForm';

const MEMORY_STORAGE = 'memory';
const REDIS_STORAGE = 'redis';
const DYNAMODB_STORAGE = 'dynamodb';
const DEFAULT_SESSION_DURATION = 300;

const SessionStorageForm = (props) => {
  const { loading, settings, error, onSubmit } = props;

  const handleSubmit = (values) => {
    onSubmit(update(values, {
      sessionDuration: value => Number(value) || DEFAULT_SESSION_DURATION,
      db: value => Number(value) || undefined,
    }));
  };

  const form = useForm({
    type: '',
    sessionDuration: '',
    address: '',
    password: '',
    db: '',
    region: '',
    endpoint: '',
  }, validateSessionStorageForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(update(settings, {
      sessionDuration: value => value === undefined ? '' : value.toString(),
      db: value => value === undefined ? '' : value.toString(),
    }));
  }, [settings]);

  return (
    <form className="iap-settings-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      <Field label="Storage Type">
        <Select
          value={form.values.type}
          disabled={loading}
          onChange={value => form.setValue('type', value)}
          placeholder="Select Storage Type"
          errorMessage={form.errors.type}
        >
          <Option value={MEMORY_STORAGE} title="Memory" />
          <Option value={REDIS_STORAGE} title="Redis" />
          <Option value={DYNAMODB_STORAGE} title="DynamoDB" />
        </Select>
      </Field>

      <Field label="Session Duration">
        <Input
          name="sessionDuration"
          value={form.values.sessionDuration}
          placeholder="Specify session duration in seconds"
          onChange={form.handleChange}
          disabled={loading}
          errorMessage={form.errors.sessionDuration}
        />
      </Field>

      {form.values.type === REDIS_STORAGE && (
        <>
          <Field label="Address">
            <Input
              name="address"
              value={form.values.address}
              placeholder="Specify address"
              onChange={form.handleChange}
              disabled={loading}
              errorMessage={form.errors.address}
            />
          </Field>
          <Field label="Password">
            <Input
              name="password"
              value={form.values.password}
              disabled={loading}
              placeholder="Specify password"
              onChange={form.handleChange}
              errorMessage={form.errors.password}
            />
          </Field>
          <Field label="DB">
            <Input
              name="db"
              value={form.values.db}
              disabled={loading}
              placeholder="Specify DB"
              onChange={form.handleChange}
              errorMessage={form.errors.db}
            />
          </Field>
        </>
      )}

      {form.values.type === DYNAMODB_STORAGE && (
        <>
          <Field label="Region">
            <Input
              name="region"
              value={form.values.region}
              placeholder="Specify region"
              disabled={loading}
              onChange={form.handleChange}
              errorMessage={form.errors.region}
            />
          </Field>
          <Field label="Endpoint" subtext="Can be figured out automatically from the region">
            <Input
              name="endpoint"
              value={form.values.endpoint}
              placeholder="Specify endpoint"
              disabled={loading}
              onChange={form.handleChange}
              errorMessage={form.errors.endpoint}
            />
          </Field>
        </>
      )}

      <footer className="iap-settings-form__footer">
        <Button
          type="submit"
          error={!loading && error}
          disabled={loading || !form.values.type}
          Icon={loading ? LoadingIcon : SaveIcon}
        >
          Save Changes
        </Button>
      </footer>
    </form>
  );
};

export default SessionStorageForm;
