import PropTypes from 'prop-types';
import React from 'react';
import Field from '~/components/shared/Field';
import Input from '~/components/shared/Input';
import { Option, Select } from '~/components/shared/Select';
import PreviewPreloader from './PreviewPreloader';

const MONGO_DB = 'mongodb';
const DYNAMO_DB = 'dynamodb';
const BOLT_DB = 'boltdb';

const createPlaceholder = fieldName => ({
  path: 'e.g. ./db.db',
  region: 'e.g. ap-northeast-3',
  endpoint: 'e.g. localhost:27017',
  connection: 'e.g. mongodb://localhost:27017',
  database: 'e.g. identifo',
}[fieldName]);

const Preview = ({ fetching, settings }) => {
  if (fetching || !settings) {
    return <PreviewPreloader />;
  }

  return (
    <div className="iap-db-connection-section">
      <div className="iap-db-form">
        <Field label="Database type">
          <Select
            name="type"
            value={settings.type}
            placeholder="Select Database Type"
            disabled
          >
            <Option value={MONGO_DB} title="Mongo DB" />
            <Option value={DYNAMO_DB} title="Dynamo DB" />
            <Option value={BOLT_DB} title="Bolt DB" />
          </Select>
        </Field>
        {Object.entries(settings[settings.type]).map(([key, value]) => (
          <Field label={key} key={key}>
            <Input
              name="region"
              value={value}
              placeholder={createPlaceholder(key)}
              disabled
            />
          </Field>
        ))}
      </div>
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
