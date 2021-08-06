import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import useDropdown from 'use-dropdown';
import Button from '~/components/shared/Button';
import DropdownIcon from '~/components/icons/DropdownIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';

const STICKING_POINT = 60;

const ActionsButton = (props) => {
  const [containerRef, isOpen, open] = useDropdown();
  const [stuck, setStuck] = useState(false);
  const [offsetX, setOffsetX] = useState(0);

  const handleScroll = () => {
    const { scrollTop } = document.documentElement;

    if (scrollTop > STICKING_POINT) {
      const { x } = containerRef.current.getBoundingClientRect();
      setOffsetX(x);
      setStuck(true);
    }

    if (scrollTop < STICKING_POINT) {
      setStuck(false);
      setOffsetX('auto');
    }
  };

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const rootClassName = classnames('iap-actions-btn', {
    'iap-actions-btn--stuck': stuck,
  });

  return (
    <div
      ref={containerRef}
      className={rootClassName}
      style={{ left: offsetX }}
    >
      <Button
        Icon={props.loading ? LoadingIcon : DropdownIcon}
        disabled={props.loading}
        iconClassName="iap-actions-btn__icon"
        onClick={open}
      >
        {props.text}
      </Button>
      {isOpen && (
        <div className="iap-actions-btn__actions">
          {props.actions.map((action) => {
            return (
              <button
                key={action.title}
                type="button"
                className="iap-actions-btn__action"
                onClick={action.onClick}
              >
                {action.title}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};

ActionsButton.propTypes = {
  actions: PropTypes.arrayOf(PropTypes.shape({
    title: PropTypes.string,
    onClick: PropTypes.func,
  })).isRequired,
  loading: PropTypes.bool,
  text: PropTypes.string,
};

ActionsButton.defaultProps = {
  loading: false,
  text: 'Actions',
};

export default ActionsButton;
