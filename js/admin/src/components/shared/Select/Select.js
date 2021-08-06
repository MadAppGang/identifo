import React, { useState, useEffect, Children } from 'react';
import PropTypes from 'prop-types';
import useDropdown from 'use-dropdown';
import DropdownIcon from '~/components/icons/DropdownIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import Input from '~/components/shared/Input';

const getDisplayValue = (value, children) => {
  const child = Children
    .toArray(children)
    .find(ch => ch.props.value === value);

  return child ? child.props.title : '';
};

let loadingTimeout;

const Select = (props) => {
  const { children, value, disabled, placeholder, errorMessage } = props;
  const [containerRef, isOpen, open, close] = useDropdown();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (loading) {
      setLoading(false);
    }

    clearTimeout(loadingTimeout);
  }, [value]);

  return (
    <div className="iap-db-dropdown" ref={containerRef}>
      <Input
        placeholder={placeholder}
        style={{ caretColor: 'transparent', cursor: 'pointer' }}
        value={getDisplayValue(value, props.children)}
        disabled={disabled}
        onFocus={open}
        errorMessage={errorMessage}
        renderButton={() => {
          if (loading) {
            return (
              <LoadingIcon className="iap-select-icon" />
            );
          }

          return (
            <DropdownIcon
              className="iap-select-icon"
              onClick={isOpen ? close : open}
            />
          );
        }}
      />

      {isOpen && (
        <div className="iap-db-dropdown__options">
          {Children.map(children, (child) => {
            const extraProps = {
              onClick: () => {
                close();

                if (child.props.value === value) {
                  return;
                }

                props.onChange(child.props.value);
                loadingTimeout = setTimeout(setLoading, 70, true);
              },
              active: value === child.props.value,
            };

            return React.cloneElement(child, extraProps);
          })}
        </div>
      )}
    </div>
  );
};

Select.propTypes = {
  disabled: PropTypes.bool,
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  errorMessage: PropTypes.string,
  placeholder: PropTypes.string,
  children: PropTypes.node.isRequired,
};

Select.defaultProps = {
  disabled: false,
  placeholder: 'Select',
  errorMessage: '',
};

export default Select;
