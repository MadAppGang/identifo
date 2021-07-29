import React, {useState} from 'react';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import LoadingIcon from '~/components/icons/LoadingIcon';
import MultipleInput from '~/components/shared/MultipleInput';
import SaveIcon from '~/components/icons/SaveIcon';
import useForm from '~/hooks/useForm';
import SecretField from './SecretField';
import { Select, Option } from '~/components/shared/Select';


const isValidUrl = string => {
  try {
    new URL(string);
  } catch (_) {
    return false;  
  }
  return true;
}

const TokenSettingsForm = ({ application, loading, onCancel, onSubmit }) => {
  
  const initialValues = {
    tokenLifespan: application.token_lifespan || '',
    refreshTokenLifespan: application.refresh_token_lifespan || '',
    inviteTokenLifespan: application.invite_token_lifespan || '',
    tokenPayload: application.token_payload || [],
    tokenPayloadService: application.token_payload_service || 'none',
    tokenPayloadServiceHttpSetting: application.token_payload_service_http_settings || {url:'', secret:''}
  };

  const [validation, setValidation] = useState({
    url: '',
  });

  const handleSubmit = (values) => {
    if ((values.tokenPayloadService == 'http') && !isValidUrl(values.tokenPayloadServiceHttpSetting.url)) {
        setValidation({...validation, url: 'Invalid URL for service, service URL is required for http service'})
        return;
    }

    onSubmit({
      ...application,
      token_lifespan: Number(values.tokenLifespan) || undefined,
      refresh_token_lifespan: Number(values.refreshTokenLifespan) || undefined,
      invite_token_lifespan: Number(values.inviteTokenLifespan) || undefined,
      token_payload: values.tokenPayload,
      token_payload_service: values.tokenPayloadService || 'none',
      token_payload_service_http_settings: values.tokenPayloadServiceHttpSetting || {}
    });
  };

  const form = useForm(initialValues, null, handleSubmit);

  React.useEffect(() => {
    if (!application) return;

    form.setValues({
      tokenLifespan: application.token_lifespan,
      refreshTokenLifespan: application.refresh_token_lifespan,
      inviteTokenLifespan: application.invite_token_lifespan,
      tokenPayload: application.token_payload,
      tokenPayloadService: application.token_payload_service,
      tokenPayloadServiceHttpSetting: application.token_payload_service_http_settings  
    });
  }, [application]);


  
  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      <Field label="Access Token Lifespan">
        <Input
          name="tokenLifespan"
          value={form.values.tokenLifespan}
          autoComplete="off"
          placeholder="Lifespan in seconds"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Refresh Token Lifespan">
        <Input
          name="refreshTokenLifespan"
          value={form.values.refreshTokenLifespan}
          autoComplete="off"
          placeholder="Lifespan in seconds"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Invite Token Lifespan">
        <Input
          name="inviteTokenLifespan"
          value={form.values.inviteTokenLifespan}
          autoComplete="off"
          placeholder="Lifespan in seconds"
          onChange={form.handleChange}
          disabled={loading}
        />
      </Field>

      <Field label="Token Payload">
        <MultipleInput
          values={form.values.tokenPayload}
          autoComplete="off"
          placeholder="Token payload"
          onChange={value => form.setValue('tokenPayload', value)}
          disabled={loading}
        />
      </Field>

      <Field label="Token Payload service">
        <Select
          name="tokenPayloadService"
          value={form.values.tokenPayloadService}
          disabled={loading}
          onChange={value => form.setValue('tokenPayloadService', value)}
          placeholder="Select Token Payload Service"
        >
          <Option value="none" title="None" />
          <Option value="http" title="External HTTP service" />
          {/* <Option value="plugin" title="Plugin" /> */}
        </Select>
      </Field>

      {form.values.tokenPayloadService === 'http' && (
        <Field label="URL">
          <Input 
            value={form.values.tokenPayloadServiceHttpSetting.url || ''}
            autoComplete="off"
            placeholder="Service URL"
            disabled={loading}
            errorMessage={validation.url}
            onValue={value => {
              form.setValue('tokenPayloadServiceHttpSetting', {...form.values.tokenPayloadServiceHttpSetting, url: value});
              setValidation({...validation, url: ''});
            }}
          />
        </Field>
      )}

      {form.values.tokenPayloadService === 'http' && (
        <SecretField 
          label="Signature secret"
          value={form.values.tokenPayloadServiceHttpSetting.secret || ''}
          onChange={value => form.setValue('tokenPayloadServiceHttpSetting', {...form.values.tokenPayloadServiceHttpSetting, secret: value})}
        />
      )}

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          disabled={loading}
          Icon={loading ? LoadingIcon : SaveIcon}
        >
          Save Changes
        </Button>
        <Button transparent disabled={loading} onClick={onCancel}>
          Cancel
        </Button>
      </footer>
    </form>
  );
};

export default TokenSettingsForm;
