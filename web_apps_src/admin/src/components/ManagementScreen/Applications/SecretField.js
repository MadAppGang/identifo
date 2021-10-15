import React, { useState } from 'react';
import PropTypes from 'prop-types';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import RevealIcon from '~/components/icons/RevealIcon';
import HideIcon from '~/components/icons/HideIcon';
import { generateSecret } from '~/utils';

const SecretField = ({ value, onChange, label = 'Client Secret' }) => {
  const [reveal, setReveal] = useState(false);


  const iconProps = {
    className: 'iap-apps-form__reveal-secret-btn',
    style: { cursor: 'pointer' },
    onClick: () => setReveal(!reveal),
  };

  const Icon = reveal ? RevealIcon : HideIcon;

  return (
    <Field
      label={label}
      Icon={<Icon {...iconProps} />}
    >
      <div className="iap-apps-form__secret-field">
        <span className="iap-apps-form__secret-value">
          {reveal ? value : '•'.repeat(value.length * 2)}
        </span>

        <Button
          disabled={!reveal}
          onClick={() => onChange(generateSecret())}
        >
          Generate
        </Button>
      </div>
    </Field>
  );
};

SecretField.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
};

export default SecretField;
