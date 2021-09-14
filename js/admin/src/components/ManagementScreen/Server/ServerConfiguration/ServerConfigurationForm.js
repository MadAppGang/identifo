import React from 'react';
import { useSelector } from 'react-redux';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import { getSettingsConfig } from '~/modules/settings/selectors';

const separeteCamelCase = (str, separator = ' ') => str.replace(/[A-Z]/gm, match => `${separator}${match}`);

const ServerConfigurationForm = () => {
  const settings = useSelector(getSettingsConfig);

  return (
    <form className="iap-apps-form">
      <Field label="Storage Type">
        <Input
          value={settings.type}
          autoComplete="off"
          placeholder="Storage type"
          disabled
        />
      </Field>

      {Object.keys(settings[settings.type]).map(key => (
        <Field key={key} label={separeteCamelCase(key)}>
          <Input
            value={settings[settings.type][key]}
            disabled
          />
        </Field>
      ))}
    </form>
  );
};

ServerConfigurationForm.defaultProps = {
  settings: {},
  loading: false,
  error: null,
};

export default ServerConfigurationForm;
