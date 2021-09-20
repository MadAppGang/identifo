import React from 'react';
import update from '@madappgang/update-by-path';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';

const MEMORY_STORAGE = 'memory';
const REDIS_STORAGE = 'redis';
const DYNAMODB_STORAGE = 'dynamo';

const validateForm = (values) => {
  const errors = {};
  if (values[REDIS_STORAGE].db && Number.isNaN(+values[REDIS_STORAGE].db)) {
    errors[REDIS_STORAGE] = { db: 'Number only' };
  }
  return errors;
};

const initialValues = {
  [REDIS_STORAGE]: { address: '', password: '', db: '' },
  [DYNAMODB_STORAGE]: { region: '', endpoint: '' },
  type: '',
};
const SessionStorageForm = (props) => {
  const { loading, settings, error, onSubmit } = props;

  const handleSubmit = (values) => {
    const payload = values.type === REDIS_STORAGE
      ? { ...values, [REDIS_STORAGE]: { ...values[REDIS_STORAGE], db: +values[REDIS_STORAGE].db } }
      : values;
    onSubmit(update(settings, payload));
  };

  // TODO: Nikita K add validation
  const form = useForm(initialValues, validateForm, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(settings);
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

      {form.values.type === REDIS_STORAGE && (
        <>
          <Field label="Address">
            <Input
              name="redis.address"
              value={form.values.redis.address}
              placeholder="Specify address"
              onChange={form.handleChange}
              disabled={loading}
              errorMessage={form.errors.address}
            />
          </Field>
          <Field label="Password">
            <Input
              name="redis.password"
              value={form.values.redis.password}
              disabled={loading}
              placeholder="Specify password"
              onChange={form.handleChange}
              errorMessage={form.errors.password}
            />
          </Field>
          <Field label="DB">
            <Input
              name="redis.db"
              value={form.values.redis.db}
              disabled={loading}
              placeholder="Specify DB"
              onChange={form.handleChange}
              errorMessage={form.errors[REDIS_STORAGE] && form.errors[REDIS_STORAGE].db}
            />
          </Field>
        </>
      )}

      {form.values.type === DYNAMODB_STORAGE && (
        <>
          <Field label="Region">
            <Input
              name="dynamo.region"
              value={form.values.dynamo.region}
              placeholder="Specify region"
              disabled={loading}
              onChange={form.handleChange}
              errorMessage={form.errors.region}
            />
          </Field>
          <Field label="Endpoint" subtext="Can be figured out automatically from the region">
            <Input
              name="dynamo.endpoint"
              value={form.values.dynamo.endpoint}
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
