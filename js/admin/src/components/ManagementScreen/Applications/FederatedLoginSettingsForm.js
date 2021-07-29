/* eslint-disable camelcase */

import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import omit from 'lodash.omit';
import classnames from 'classnames';
import update from '@madappgang/update-by-path';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Toggle from '~/components/shared/Toggle';
import Button from '~/components/shared/Button';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import WarningIcon from '~/components/icons/WarningIcon.svg';
import MultipleInput from '~/components/shared/MultipleInput';
import useNotifications from '~/hooks/useNotifications';

const extractValue = fn => e => fn(e.target.value);

const FederatedLoginSettingsForm = (props) => {
  const { loading, onSubmit, onCancel, federatedProviders } = props;
  const application = props.application || {};

  const [federatedLoginSettings, setFederatedLoginSettings] = useState({});

  const { notifyFailure } = useNotifications();

  useEffect(() => {
    setFederatedLoginSettings(application.federated_login_settings || {});
  }, [props.application]);

  const handleParamsChange = (provider, field, value) => {
    const params = { ...federatedLoginSettings[provider].params, [field]: value };

    setFederatedLoginSettings({
      ...federatedLoginSettings,
      [provider]: { ...federatedLoginSettings[provider], params },
    });
  };

  const handleScopesChange = (provider, value) => {
    setFederatedLoginSettings({
      ...federatedLoginSettings,
      [provider]: { ...federatedLoginSettings[provider], scopes: [...value] },
    });
  };

  // Each param contain comma separated options.
  // For example when param set to "PKCS8PrivateKey,textarea,optional"
  // need to use textarea control and this param is optional.
  // Options is unsorted and can be in any order.
  // First element of comma separated list is always field name.
  const extractParam = (param) => {
    const params = param.split(',');
    return {
      fieldName: params[0],
      textarea: !!params.find(p => p === 'textarea'),
      optional: !!params.find(p => p === 'optional'),
    };
  };

  const toggleProvider = (value) => {
    if (federatedLoginSettings[value]) {
      setFederatedLoginSettings(omit(federatedLoginSettings, value));
      return;
    }

    // create params fields for current provider
    const initialParams = Object.fromEntries(federatedProviders[value].params.map(p => [extractParam(p).fieldName, '']));

    // create initial structure for settings object
    setFederatedLoginSettings({
      ...federatedLoginSettings,
      [value]: {
        params: { ...initialParams },
        scopes: [],
      },
    });
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const isSettingsProvided = Object.values(federatedLoginSettings).findIndex((v) => {
      return Object.values(v.params).includes('');
    }) < 0;

    if (!isSettingsProvided) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'All arguments must be provided!',
      });
      return;
    }

    onSubmit(update(application, {
      federated_login_settings: federatedLoginSettings,
    }));
  };

  return (
    <form className="iap-apps-form" onSubmit={handleSubmit}>
      <div className="iap-apps-form__note">
        <WarningIcon className="iap-apps-form__note-icon" />
        <p>
          Note that these settings take effect only when federated login is enabled in
          <Link className="iap-apps-form__note-link" to="/management/settings">
            Login Types
          </Link>
          settings.
        </p>
      </div>

      {Object.entries(federatedProviders).map((provider) => {
        const currentProvider = provider[0];

        const isActive = currentProvider in federatedLoginSettings;

        const providerClassName = classnames({
          'iap-apps-form__provider': true,
          'iap-apps-form__provider--open': isActive,
        });

        return (
          <div key={currentProvider} className={providerClassName}>
            <Toggle
              label={provider[1].string}
              value={isActive}
              onChange={() => toggleProvider(currentProvider)}
            />

            {isActive && (
              <>
                {federatedProviders[currentProvider].params.map((param) => {
                  const { fieldName, textarea } = extractParam(param);
                  return (
                    <Field label={fieldName} key={fieldName}>
                      {!textarea && (
                        <Input
                          value={federatedLoginSettings[currentProvider].params[fieldName]}
                          autoComplete="off"
                          placeholder={`Enter ${fieldName}`}
                          onChange={extractValue(v => handleParamsChange(currentProvider, fieldName, v))}
                        />
                      )}
                      {textarea && (
                        <textarea
                          value={federatedLoginSettings[currentProvider].params[fieldName]}
                          autoComplete="off"
                          placeholder={`Enter ${fieldName}`}
                          className="iap-login-form__input--textarea iap-login-form__input"
                          onChange={extractValue(v => handleParamsChange(currentProvider, fieldName, v))}
                        />
                      )}
                    </Field>
                  );
                })}

                <Field label="Scopes">
                  <MultipleInput
                    values={federatedLoginSettings[currentProvider].scopes}
                    placeholder="Hit Enter to add scope"
                    onChange={v => handleScopesChange(currentProvider, v)}
                  />
                </Field>

                <div className="iap-apps-form__note">
                  <WarningIcon className="iap-apps-form__note-icon" />
                  <p>
                    Don&#39;t forget to add redirect URI
                    {` ${window.location.origin}/web/login?appId=${application.id}&provider=${currentProvider} `}
                    to auth provider settings.
                  </p>
                </div>
              </>
            )
            }
          </div>
        );
      })}

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

export default FederatedLoginSettingsForm;
