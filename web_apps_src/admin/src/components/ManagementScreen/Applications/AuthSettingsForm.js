import React, { useState, useEffect } from 'react';
import update from '@madappgang/update-by-path';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import { Select, Option } from '~/components/shared/Select';
import MultipleInput from '~/components/shared/MultipleInput';
import CasbinEditor from './CasbinEditor';

const extractValue = fn => event => fn(event.target.value);

const ApplicationAuthSettings = (props) => {
  const { loading, onSubmit, onCancel } = props;

  const application = props.application || {};

  const [authWay, setAuthWay] = useState(application.authorization_way || '');
  const [authModel, setAuthModel] = useState(application.authorization_model || '');
  const [authPolicy, setAuthPolicy] = useState(application.authorization_policy || '');
  const [defaultRole, setDefaultRole] = useState(application.new_user_default_role || '');
  const [whitelist, setWhitelist] = useState(application.roles_whitelist || []);
  const [blacklist, setBlacklist] = useState(application.roles_blacklist || []);

  useEffect(() => {
    setAuthWay(application.authorization_way || '');
    setAuthModel(application.authorization_model || '');
    setAuthPolicy(application.authorization_policy || '');
    setDefaultRole(application.new_user_default_role || '');
    setWhitelist(application.roles_whitelist || []);
    setBlacklist(application.roles_blacklist || []);
  }, [props.application]);

  const handleSubmit = (event) => {
    event.preventDefault();

    onSubmit(update(application, {
      authorization_way: authWay,
      authorization_model: authModel,
      authorization_policy: authPolicy,
      new_user_default_role: defaultRole,
      roles_whitelist: whitelist,
    }));
  };

  return (
    <form className="iap-apps-form" onSubmit={handleSubmit}>
      <Field label="Authorization Way">
        <Select
          value={authWay}
          disabled={loading}
          onChange={setAuthWay}
          placeholder="Select Authorization Way"
        >
          <Option value="no_authorization" title="No Authorization" />
          <Option value="internal" title="Internal" />
          <Option value="whitelist" title="Whitelist" />
          <Option value="blacklist" title="Blacklist" />
          <Option value="external" title="External" />
        </Select>
      </Field>

      {authWay === 'internal' && (
        <Field label="Casbin Model">
          <CasbinEditor value={authModel} onChange={setAuthModel} />
        </Field>
      )}

      {authWay === 'internal' && (
        <Field label="Casbin Policy">
          <CasbinEditor value={authPolicy} onChange={setAuthPolicy} />
        </Field>
      )}

      {authWay === 'whitelist' && (
        <Field label="Roles Whitelist">
          <MultipleInput
            values={whitelist}
            placeholder="Hit Enter to add role"
            onChange={setWhitelist}
          />
        </Field>
      )}

      {authWay === 'blacklist' && (
        <Field label="Roles Blacklist">
          <MultipleInput
            values={blacklist}
            placeholder="Hit Enter to add role"
            onChange={setBlacklist}
          />
        </Field>
      )}

      <Field label="New User Default Role">
        <Input
          value={defaultRole}
          autoComplete="off"
          placeholder="User role"
          onChange={extractValue(setDefaultRole)}
        />
      </Field>

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

export default ApplicationAuthSettings;
