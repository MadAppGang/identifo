import React from 'react';
import Field from '~/components/shared/Field';
import MultipleInput from '~/components/shared/MultipleInput';
import useForm from '~/hooks/useForm';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import Toggle from '~/components/shared/Toggle';
import { useSelector } from 'react-redux';
import update from '@madappgang/update-by-path';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import { toDeepCase } from '~/utils/apiMapper';

const initalValues = {
  newUserDefaultScopes: [],
  newUserDefaultRole: '',
  registrationForbidden: false,
  anonymousRegistrationAllowed: false,
};
export const RegistrationSettingsForm = (props) => {
  const { loading, error, onSubmit, onCancel } = props;
  const application = useSelector(s => s.selectedApplication.application);
  const submitHandler = (v) => {
    onSubmit(update(application, toDeepCase(v, 'snake')));
  };
  const form = useForm(initalValues, null, submitHandler);
  React.useEffect(() => {
    if (application) {
      form.setValues(update(form.values, {
        registrationForbidden: application.registration_forbidden || false,
        anonymousRegistrationAllowed: !!application.anonymous_registration_allowed,
        newUserDefaultScopes: application.new_user_default_scopes || [],
        newUserDefaultRole: application.new_user_default_role || '',
      }));
    }
  }, [application]);
  return (
    <form className="iap-apps-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}
      <Field label="New user default scopes">
        <MultipleInput
          values={form.values.newUserDefaultScopes}
          placeholder="Hit Enter to add scope"
          onChange={v => form.setValue('newUserDefaultScopes', v)}
        />
      </Field>

      <Field label="New user default role">
        <Input
          value={form.values.newUserDefaultRole}
          name="newUserDefaultRole"
          autoComplete="off"
          placeholder="Enter new user default role"
          onChange={form.handleChange}
        />
      </Field>

      <Toggle
        label="Allow Registration"
        value={!form.values.registrationForbidden}
        onChange={v => form.setValue('registrationForbidden', !v)}
      />

      <Toggle
        label="Allow Anonymous Registration"
        value={form.values.anonymousRegistrationAllowed}
        onChange={v => form.setValue('anonymousRegistrationAllowed', v)}
      />

      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          disabled={loading}
          Icon={loading ? LoadingIcon : SaveIcon}
        >
          Save Changes
        </Button>
        <Button
          transparent
          disabled={loading}
          onClick={onCancel}
        >
          Cancel
        </Button>
      </footer>
    </form>
  );
};
