import copy from 'copy-to-clipboard';
import React, { useEffect, useState } from 'react';
import CopyIcon from '~/components/icons/CopyIcon';
import EditIcon from '~/components/icons/EditIcon';
import KeyIcon from '~/components/icons/Key';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import PaswordIcon from '~/components/icons/ShowPassword';
import Button from '~/components/shared/Button';
import Field from '~/components/shared/Field';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import useForm from '~/hooks/useForm';

const validationSchema = (values) => {
  const errors = {};
  if (!values.private) {
    errors.private = 'Private key couldn`t be empty';
  }
  return errors;
};

const initialValues = {
  private: '',
  public: '',
  alg: '',
};

const FormFooter = ({ changing, loading, error, onGenerateKey, onEdit, onCancel }) => {
  if (!changing) {
    return (
      <Button
        onClick={onEdit}
        Icon={loading ? LoadingIcon : EditIcon}
        disabled={loading}
        error={!loading && !!error}
        key="edit"
      >
        Edit
      </Button>
    );
  }
  return (
    <>
      <Button
        type="submit"
        Icon={loading ? LoadingIcon : SaveIcon}
        disabled={loading}
        error={!loading && !!error}
        key="submit"
      >
        Save Changes
      </Button>
      <Button
        type="button"
        onClick={onGenerateKey}
        Icon={loading ? LoadingIcon : KeyIcon}
        disabled={loading}
        key="generate"
        error
      >
        Generate Key
      </Button>
      <Button
        type="button"
        onClick={onCancel}
        disabled={loading}
        key="cancel"
        transparent
      >
        Cancel
      </Button>
    </>
  );
};

export const JWTSettingsForm = ({
  error, loading, settings, onShowPassword,
  onGenerateKey, onSubmit,
}) => {
  const form = useForm(initialValues, validationSchema, onSubmit);
  const [changing, setChanging] = useState(false);

  const showPasswordHandler = async () => {
    if (!settings || !settings.private) {
      await onShowPassword();
    } else if (!form.values.private) {
      form.setValues(settings);
    } else {
      form.setValues({ ...settings, private: '' });
    }
  };

  const cancelHandelr = () => {
    setChanging(false);
    form.setValues({ ...settings, private: '' });
  };

  useEffect(() => {
    if (settings) {
      form.setValues(settings);
    }
  }, [settings]);

  return (
    <form className="iap-apps-form iap-jwt-settings-form" onSubmit={form.handleSubmit}>
      {!!error && (
        <FormErrorMessage error={error} />
      )}
      <Field label="Private key" subtext="Please paste pkcs8 pem private key">
        <div className="iap-jwt-settings-form--field">
          <Input
            className="iap-login-form__input iap-login-form__input--textarea"
            name="private"
            as={form.values.private ? 'textarea' : 'input'}
            placeholder="Click reveal to show private key"
            value={form.values.private}
            onChange={form.handleChange}
            autoComplete="off"
            errorMessage={form.errors.private}
            disabled={!form.values.private || !changing}
          />
          <PaswordIcon className="iap-jwt-settings-form-action-btn" onClick={showPasswordHandler} />
        </div>
      </Field>
      <Field label="Public key">
        <div className="iap-jwt-settings-form--field">
          <textarea
            className="iap-login-form__input iap-login-form__input--textarea"
            name="public"
            placeholder="Enter your public key"
            value={form.values.public}
            onChange={form.handleChange}
            autoComplete="off"
            disabled
          />
          <CopyIcon className="iap-jwt-settings-form-action-btn" onClick={() => copy(settings.public)} />
        </div>
      </Field>
      <Field label="Algorithm">
        <Input
          name="private"
          placeholder="Selected algorithm"
          value={form.values.alg}
          onChange={form.handleChange}
          autoComplete="off"
          disabled
        />
      </Field>
      <footer className="iap-apps-form__footer">
        <FormFooter
          changing={changing}
          onGenerateKey={onGenerateKey}
          onEdit={() => setChanging(true)}
          onCancel={cancelHandelr}
          loading={loading}
          error={error}
        />
      </footer>
    </form>
  );
};
