import React from 'react';
import PropTypes from 'prop-types';
import PreviewPreloader from './PreviewPreloader';
import PreviewField from '~/components/shared/PreviewField';

const MONGO_DB = 'mongodb';
const DYNAMO_DB = 'dynamodb';
const BOLT_DB = 'boltdb';

const displayDatabaseType = type => ({
  [MONGO_DB]: 'MongoDB',
  [DYNAMO_DB]: 'DynamoDB',
  [BOLT_DB]: 'BoltDB',
}[type]);

const Preview = ({ fetching, settings }) => {
  if (fetching || !settings) {
    return <PreviewPreloader />;
  }

  const { type } = settings;

  return (
    <div className="iap-section__info">
      <PreviewField
        label="Database Type"
        value={displayDatabaseType(settings.type)}
      />

      {type === MONGO_DB && (
        <PreviewField
          label="Database Name"
          value={settings.name}
        />
      )}

      {type === DYNAMO_DB && (
        <PreviewField
          label="Region"
          value={settings.region}
        />
      )}

      {type !== BOLT_DB && (
        <PreviewField
          label="Endpoint"
          value={settings.endpoint}
        />
      )}

      {type === BOLT_DB && (
        <PreviewField
          label="Path"
          value={settings.path}
        />
      )}
    </div>
  );
};

Preview.propTypes = {
  fetching: PropTypes.bool,
  settings: PropTypes.shape({
    endpoint: PropTypes.string,
    name: PropTypes.string,
    type: PropTypes.string,
  }),
};

Preview.defaultProps = {
  fetching: false,
  settings: null,
};

export default Preview;
