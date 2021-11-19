import update from "@madappgang/update-by-path";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import Field from "~/components/shared/Field";
import Input from "~/components/shared/Input";
import { Option, Select } from "~/components/shared/Select";
import useProgressBar from "~/hooks/useProgressBar";
import { updateServerSettings } from "~/modules/settings/actions";
import { getLoginSettings } from "~/modules/settings/selectors";
import LoginTypesTable from "./LoginTypesTable";

const LoginTypesSection = () => {
  const dispatch = useDispatch();
  const settings = useSelector(getLoginSettings);
  const { setProgress } = useProgressBar();

  useEffect(() => {
    // TODO: Nikita K removee this uef
    setProgress(100);
  }, []);

  const handleChange = (type, enabled) => {
    const nextSettings = update(settings, `loginWith.${type}`, enabled);
    dispatch(updateServerSettings({ login: nextSettings }));
  };

  const handleTfaTypeChange = (value) => {
    const nextSettings = update(settings, { tfaType: value });
    dispatch(updateServerSettings({ login: nextSettings }));
  };

  const handleInput = ({ target }) => {
    const nextSettings = update(settings, {
      [target.name]: Number(target.value),
    });
    dispatch(updateServerSettings({ login: nextSettings }));
  };

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">Login Types</p>

      <p className="iap-management-section__description">
        These settings allow to turn off undesirable login endpoints.
      </p>

      <div className="iap-settings-section">
        <div className="section-field">
          <Field label="2FA Type">
            <Select
              value={settings.tfaType}
              onChange={handleTfaTypeChange}
              placeholder="Select 2FA Type"
            >
              <Option value="app" title="App" />
              <Option value="sms" title="SMS" />
              <Option value="email" title="Email" />
            </Select>
          </Field>
          <Field label="2FA Resend timeout (seconds)">
            <Input
              name="tfaResendTimeout"
              value={settings.tfaResendTimeout}
              autoComplete="off"
              placeholder="Timeout to show send again"
              onChange={handleInput}
            />
          </Field>
        </div>
        <LoginTypesTable types={settings.loginWith} onChange={handleChange} />
      </div>
    </section>
  );
};

export default LoginTypesSection;
