/* eslint-disable max-len */
import copy from 'copy-to-clipboard';
import React, { useEffect } from 'react';
import CopyIcon from '~/components/icons/CopyIcon';
import KeyIcon from '~/components/icons/Key';
import LoadingIcon from '~/components/icons/LoadingIcon';
import SaveIcon from '~/components/icons/SaveIcon';
import PaswordIcon from '~/components/icons/ShowPassword';
import Button from '~/components/shared/Button';
import { CollapseItem, CollapseLinks } from '~/components/shared/CollapseLink/CollapseLink';
import Field from '~/components/shared/Field';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import Input from '~/components/shared/Input';
import { Option, Select } from '~/components/shared/Select';
import useForm from '~/hooks/useForm';
import { domChangeEvent } from '~/utils';

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
export const JWTSettingsForm = ({
  error, loading, settings, onShowPassword,
  onGenerateKey, onSubmit,
}) => {
  const form = useForm(initialValues, validationSchema, onSubmit);

  const showPasswordHandler = async () => {
    if (!settings.private) {
      await onShowPassword();
    } else if (!form.values.private) {
      form.setValues(settings);
    } else {
      form.setValues({ ...settings, private: '' });
    }
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
            as="textarea"
            placeholder="Click reveal to show private key"
            value={form.values.private}
            onChange={form.handleChange}
            autoComplete="off"
            errorMessage={form.errors.private}
            disabled={!form.values.private}
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
        <Select
          value={form.values.alg}
          disabled={loading}
          onChange={v => form.handleChange(domChangeEvent('alg', v))}
          placeholder="Select Algorithm"
        >
          <Option value="RS256" title="RS256" />
          <Option value="ES256" title="ES256" />
        </Select>
      </Field>
      <CollapseLinks accordion activeTitle>
        <CollapseItem title="How to generate RS256 private key (widely supported by all framwroks)">
          <p>Paste it in your terminal:</p>
          <p>ssh-keygen -t rsa -b 2048 -m PEM -f private.pem -C &quot;identifo@madappgang.com&quot; -N &quot;&quot;</p>
          <p>rm private.pem.pub</p>
          <p>openssl rsa -in private.pem -pubout -outform PEM -out public.pem</p>
        </CollapseItem>
        <CollapseItem title="How to generate EC secp256k1  private key (RECOMMENDED)">
          <p>Paste it in your terminal:</p>
          <p>openssl ecparam -name prime256v1 -genkey -noout -out private_ec.pem</p>
          <p>openssl pkcs8 -topk8 -nocrypt -inform PEM -outform PEM -in private_ec.pem -out private.pem</p>
          <p>rm private_ec.pem</p>
          <p>openssl ec -in private.pem -pubout -out public.pem</p>
        </CollapseItem>
      </CollapseLinks>
      <footer className="iap-apps-form__footer">
        <Button
          type="submit"
          Icon={loading ? LoadingIcon : SaveIcon}
          disabled={loading}
          error={!loading && !!error}
        >
          Save Changes
        </Button>
        <Button
          type="button"
          onClick={() => onGenerateKey(form.values.alg)}
          Icon={loading ? LoadingIcon : KeyIcon}
          disabled={loading}
          error
        >
          Generate Key
        </Button>
      </footer>
    </form>
  );
};
