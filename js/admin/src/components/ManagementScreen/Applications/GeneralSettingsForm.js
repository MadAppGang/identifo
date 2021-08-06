import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import * as Validation from '@dprovodnikov/validation';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import validationRules from './validationRules';
import { Select, Option } from '~/components/shared/Select';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Toggle from '~/components/shared/Toggle';
import SecretField from './SecretField';
import MultipleInput from '~/components/shared/MultipleInput';

const extractValue = fn => event => fn(event.target.value);
const validate = Validation.applyRules(validationRules);

const ApplicationGeneralSettingsForm = (props) => {
  const { loading, error, excludeFields, onSubmit, onCancel } = props;
  const application = props.application || {};

  const [redirectUrls, setRedirectUrls] = useState(application.redirect_url || []);
  const [offline, setOffline] = useState(application.offline || false);
  const [type, setType] = useState(application.type || 'web');
  const [name, setName] = useState(application.name || '');
  const [description, setDescription] = useState(application.description || '');
  const [secret, setSecret] = useState(application.secret || '');
  const [allowRegistration, setAllowRegistration] = useState(
    !application.registration_forbidden || false,
  );
  const [allowAnonymousRegistration, setAllowAnonymousRegistration] = useState(
    application.anonymous_registration_allowed || false,
  );
  const [tfaStatus, setTfaStatus] = useState(application.tfa_status || 'disabled');
  const [active, setActive] = useState(application.active || false);
  const [debugTfaCode, setDebugTfaCode] = useState(application.debug_tfa_code || '');
  const [scopes, setScopes] = useState(application.scopes || []);
  const [federatedLoginSettings, setFederatedLoginSettings] = useState(application.federated_login_settings || {});

  const [validation, setValidation] = useState({
    type: '',
    name: '',
    redirectUrls: '',
  });

  /* update field values after props update */
  useEffect(() => {
    if (!application) return;
    if (application.redirect_urls) setRedirectUrls(application.redirect_urls || []);
    if (application.offline) setOffline(application.offline);
    if (application.type) setType(application.type);
    if (application.name) setName(application.name);
    if (application.description) setDescription(application.description);
    if (application.secret) setSecret(application.secret);
    if (application.tfa_status) setTfaStatus(application.tfa_status);
    if (application.active) setActive(application.active);
    if (application.debug_tfa_code) setDebugTfaCode(application.debug_tfa_code);
    if (application.scopes) setScopes(application.scopes);
    if (application.federated_login_settings) setFederatedLoginSettings(application.federated_login_settings);
    setAllowRegistration(!application.registration_forbidden);
    setAllowAnonymousRegistration(application.anonymous_registration_allowed);
  }, [props.application]);

  const isExcluded = field => excludeFields.includes(field);

  const handleInput = (field, value, setValue) => {
    if (field in validation) {
      setValidation(update(validation, { [field]: '' }));
    }
    setValue(value);
  };

  const handleBlur = (field, value) => {
    const validationMessage = validate(field, { [field]: value });

    setValidation(update(validation, {
      [field]: validationMessage,
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const report = validate('all', { name, type, redirectUrls });

    if (Validation.hasError(report)) {
      setValidation(report);
      return;
    }

    onSubmit({
      offline,
      type,
      name,
      scopes,
      secret,
      active,
      description,
      tfa_status: tfaStatus,
      redirect_urls: redirectUrls,
      registration_forbidden: !allowRegistration,
      anonymous_registration_allowed: allowAnonymousRegistration,
      debug_tfa_code: debugTfaCode || undefined,
      federated_login_settings: federatedLoginSettings,
    });
  };

  return (
    <form className="iap-apps-form" onSubmit={handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}

      {!isExcluded('name') && (
        <Field label="Name">
          <Input
            value={name}
            autoComplete="off"
            placeholder="Enter name"
            onChange={extractValue(v => handleInput('name', v, setName))}
            onBlur={extractValue(v => handleBlur('name', v))}
            errorMessage={validation.name}
            disabled={loading}
          />
        </Field>
      )}

      {!isExcluded('type') && (
        <Field label="Type">
          <Select
            name="type"
            value={type}
            disabled={loading}
            onChange={setType}
            placeholder="Select Application Type"
            errorMessage={validation.type}
          >
            <Option value="web" title="Single Page Application (Web)" />
            <Option value="android" title="Android Client (Mobile)" />
            <Option value="ios" title="iOS Client (Mobile)" />
            <Option value="desktop" title="Desktop Client (Desktop)" />
          </Select>
        </Field>
      )}

      {!isExcluded('description') && (
        <Field label="Description">
          <Input
            value={description}
            autoComplete="off"
            placeholder="Enter Description"
            onChange={extractValue(v => handleInput('description', v, setDescription))}
            onBlur={extractValue(v => handleBlur('description', v))}
            disabled={loading}
          />
        </Field>
      )}

      {!isExcluded('scopes') && (
        <Field label="Scopes">
          <MultipleInput
            values={scopes}
            placeholder="Hit Enter to add scope"
            onChange={setScopes}
          />
        </Field>
      )}

      {!isExcluded('tfaStatus') && (
        <Field label="2FA Status">
          <Select
            value={tfaStatus}
            disabled={loading}
            onChange={setTfaStatus}
            placeholder="Select TFA Status"
          >
            <Option value="disabled" title="Disabled" />
            <Option value="mandaroty" title="Mandatory" />
            <Option value="optional" title="Optional" />
          </Select>
        </Field>
      )}

      {!isExcluded('secret') && (
        <SecretField value={secret} onChange={setSecret} />
      )}

      {!isExcluded('redirectUrl') && (
        <Field label="Redirect URLs">
          <MultipleInput
            values={redirectUrls}
            placeholder="Hit Enter to add url"
            onChange={v => handleInput('redirectUrls', v, setRedirectUrls)}
            errorMessage={validation.redirectUrls}
          />
        </Field>
      )}

      {!isExcluded('debugTfaCode') && (
        <Field label="Debug TFA Code">
          <Input
            value={debugTfaCode}
            autoComplete="off"
            placeholder="Debug TFA Code"
            onChange={extractValue(v => handleInput('debugTfaCode', v, setDebugTfaCode))}
            onBlur={extractValue(v => handleBlur('debugTfaCode', v))}
            disabled={loading}
          />
        </Field>
      )}

      <div>
        {!isExcluded('allowRegistration') && (
          <Toggle
            label="Allow Registration"
            value={!!allowRegistration}
            onChange={setAllowRegistration}
          />
        )}

        {!isExcluded('allowAnonymousRegistration') && (
          <Toggle
            label="Allow Anonymous Registration"
            value={!!allowAnonymousRegistration}
            onChange={setAllowAnonymousRegistration}
          />
        )}

        {!isExcluded('offline') && (
          <Toggle label="Allow Offline" value={!!offline} onChange={setOffline} />
        )}

        {!isExcluded('active') && (
          <Toggle label="Active" value={!!active} onChange={setActive} />
        )}

      </div>

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          Icon={loading ? LoadingIcon : SaveIcon}
          disabled={loading || Validation.hasError(validation)}
          error={!loading && !!error}
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

ApplicationGeneralSettingsForm.defaultProps = {
  excludeFields: [],
};

export default ApplicationGeneralSettingsForm;
