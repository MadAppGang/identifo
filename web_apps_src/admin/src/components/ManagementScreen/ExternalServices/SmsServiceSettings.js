import React from 'react';
import update from '@madappgang/update-by-path';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import Input from '~/components/shared/Input';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import { Select, Option } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';

const [TWILIO, MOCK, NEXMO, RMOBILE] = ['twilio', 'mock', 'nexmo', 'routemobile'];

const SmsServiceSettings = ({ settings, loading, onSubmit }) => {
  const initialValues = {
    type: settings ? settings.type : '',
    [TWILIO]: { accountSid: '', authToken: '', serviceSid: '' },
    [NEXMO]: { apiKey: '', apiSecret: '' },
    [RMOBILE]: { username: '', password: '', source: '', region: '' },
  };

  const handleSubmit = (values) => {
    onSubmit(update(settings, values));
  };
  // TODO:Nikita K implement form validation
  const form = useForm(initialValues, null, handleSubmit);

  React.useEffect(() => {
    if (!settings) return;

    form.setValues(settings);
  }, [settings]);

  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      <Field label="SMS Service">
        <Select
          value={form.values.type}
          disabled={loading}
          onChange={value => form.setValue('type', value)}
          placeholder="Select Service"
          errorMessage={form.errors.type}
        >
          <Option value={TWILIO} title="Twilio" />
          <Option value={NEXMO} title="Nexmo" />
          <Option value={RMOBILE} title="Route Mobile" />
          <Option value={MOCK} title="Mock" />
        </Select>
      </Field>

      {form.values.type === NEXMO && (
        <>
          <Field label="Api Key">
            <Input
              name="nexmo.apiKey"
              value={form.values.nexmo.apiKey}
              autoComplete="off"
              placeholder="Specify Nexmo api key"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.authKey}
            />
          </Field>

          <Field label="Api Secret">
            <Input
              name="nexmo.apiSecret"
              value={form.values.nexmo.apiSecret}
              autoComplete="off"
              placeholder="Specify Nexmo api secret"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.apiSecret}
            />
          </Field>
        </>
      )}

      {form.values.type === TWILIO && (
        <>
          <Field label="Auth Token">
            <Input
              name="twilio.authToken"
              value={form.values.twilio.authToken}
              autoComplete="off"
              placeholder="Specify Twilio Auth Token"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.authToken}
            />
          </Field>

          <Field label="Account SID">
            <Input
              name="twilio.accountSid"
              value={form.values.twilio.accountSid}
              autoComplete="off"
              placeholder="Specify Twilio Account SID"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.accountSid}
            />
          </Field>

          <Field label="Service SID">
            <Input
              name="twilio.serviceSid"
              value={form.values.twilio.serviceSid}
              autoComplete="off"
              placeholder="Specify Twilio Service SID"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.serviceSid}
            />
          </Field>
        </>
      )}

      {form.values.type === RMOBILE && (
        <>
          <Field label="Username">
            <Input
              name="routemobile.username"
              value={form.values.routemobile.username}
              autoComplete="off"
              placeholder="Specify Route mobile username"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.serviceSid}
            />
          </Field>

          <Field label="Password">
            <Input
              name="routemobile.password"
              value={form.values.routemobile.password}
              autoComplete="off"
              placeholder="Enter password"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.serviceSid}
            />
          </Field>

          <Field label="Source">
            <Input
              name="routemobile.source"
              value={form.values.routemobile.source}
              autoComplete="off"
              placeholder="Specify Route mobile source"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.serviceSid}
            />
          </Field>

          <Field label="Region">
            <Input
              name="routemobile.region"
              value={form.values.routemobile.region}
              autoComplete="off"
              placeholder="Specify Route mobile region"
              onChange={form.handleChange}
              onBlur={form.handleBlur}
              disabled={loading}
              errorMessage={form.errors.serviceSid}
            />
          </Field>
        </>
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

export default SmsServiceSettings;
