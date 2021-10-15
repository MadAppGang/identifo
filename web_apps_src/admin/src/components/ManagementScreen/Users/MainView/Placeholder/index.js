import React from 'react';
import PropTypes from 'prop-types';
import Button from '~/components/shared/Button';
import AddIcon from '~/components/icons/AddIcon';
import UsersIcon from '~/components/icons/UsersIcon';
import './UsersPlaceholder.css';

const UsersPlaceholder = (props) => {
  return (
    <div className="iap-section-placeholder">
      <h2 className="iap-section-placeholder__title">
        Users
      </h2>

      <UsersIcon className="iap-section-placeholder__icon" />

      <p className="iap-section-placeholder__msg">
        No users have been added to your applications.
      </p>

      <Button Icon={AddIcon} onClick={props.onCreateNewUserClick}>
        Create your first user
      </Button>
    </div>
  );
};

UsersPlaceholder.propTypes = {
  onCreateNewUserClick: PropTypes.func,
};

UsersPlaceholder.defaultProps = {
  onCreateNewUserClick: () => {},
};

export default UsersPlaceholder;
