import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import {
  CONNECTION_ESTABLISHED,
  CONNECTION_FAILED,
  CONNECTION_TEST_REQUIRED,
} from '~/modules/database/connectionReducer';
import { testConnection } from '~/modules/database/actions';
import Button from '~/components/shared/Button';
import ConnectionIcon from '~/components/icons/ConnectionIcon';
import ErrorIcon from '~/components/icons/ErrorIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';

const ConnectionState = ({ loading, checking, state, ...props }) => {
  return (
    <>
      {state === CONNECTION_TEST_REQUIRED && (
        <Button
          Icon={checking || loading ? LoadingIcon : ConnectionIcon}
          disabled={loading || checking}
          onClick={props.testConnection}
        >
          Test connection
        </Button>
      )}

      {state === CONNECTION_ESTABLISHED && (
        <div className="iap-db__connection iap-db__connection--established">
          <ConnectionIcon className="iap-db__connection-icon" />
          Connection Established
        </div>
      )}

      {state === CONNECTION_FAILED && (
        <div className="iap-db__connection iap-db__connection--failed">
          <ErrorIcon className="iap-db__connection-icon" />
          Connection Failed
        </div>
      )}
    </>
  );
};

ConnectionState.propTypes = {
  state: PropTypes.oneOf([
    CONNECTION_ESTABLISHED, CONNECTION_FAILED, CONNECTION_TEST_REQUIRED,
  ]).isRequired,
  checking: PropTypes.bool.isRequired,
  loading: PropTypes.bool,
  testConnection: PropTypes.func.isRequired,
};

ConnectionState.defaultProps = {
  loading: false,
};

const mapStateToProps = state => ({
  state: state.database.connection.state,
  checking: state.database.connection.checking,
});

const actions = {
  testConnection,
};

export { ConnectionState };

export default connect(mapStateToProps, actions)(ConnectionState);
