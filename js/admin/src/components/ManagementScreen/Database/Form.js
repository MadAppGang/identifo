import React, { useState } from 'react';
import PropTypes from 'prop-types';
import * as Validation from '@dprovodnikov/validation';
import Input from '~/components/shared/Input';
import Field from '~/components/shared/Field';
import Button from '~/components/shared/Button';
import SaveIcon from '~/components/icons/SaveIcon';
import LoadingIcon from '~/components/icons/LoadingIcon';
import databaseFormValidationRules from './validationRules';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import { Select, Option } from '~/components/shared/Select';

const MONGO_DB = 'mongodb';
const DYNAMO_DB = 'dynamodb';
const BOLT_DB = 'boltdb';

const ConnectionSettingsForm = (props) => {
  const { posting, error, settings, onCancel, onSubmit } = props;
  const { type } = settings;

  const [dbType, setDbType] = useState(type);

  const [validation, setValidation] = useState({
    region: '',
    endpoint: '',
    path: '',
    database: '',
    connection: '',
  });

  const [dbSettings, setDbSettings] = useState({
    region: type === DYNAMO_DB ? settings[type].region : '',
    endpoint: type === DYNAMO_DB ? settings[type].endpoint : '',
    path: type === BOLT_DB ? settings[type].path : '',
    database: type === MONGO_DB ? settings.mongo.database : '',
    connection: type === MONGO_DB ? settings.mongo.connection : '',
  });

  const validate = Validation.applyRules(databaseFormValidationRules);

  const changeDbType = (value) => {
    setDbType(value);
    Validation.reset(validation);
  };

  const handleInput = ({ target }) => {
    setDbSettings({ ...dbSettings, [target.name]: target.value });
    setValidation({ ...validation, [target.name]: '' });
  };

  const getFieldsToOmitDuringValidation = () => {
    switch (dbType) {
      case DYNAMO_DB: return ['name', 'path', 'connection', 'database'];
      case MONGO_DB: return ['region', 'path', 'endpoint', 'name'];
      case BOLT_DB: return ['name', 'region', 'endpoint', 'connection', 'database'];
      default: return [];
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const payload = () => {
      switch (dbType) {
        case DYNAMO_DB:
          return { region: dbSettings.region, endpoint: dbSettings.endpoint };
        case MONGO_DB:
          return { database: dbSettings.database, connection: dbSettings.connection };
        case BOLT_DB:
          return { path: dbSettings.path };
        default:
          return {};
      }
    };

    const validationReport = validate('all', dbSettings, { omit: getFieldsToOmitDuringValidation() });

    if (Validation.hasError(validationReport)) {
      setValidation(validationReport);
      return;
    }

    onSubmit(
      dbType === MONGO_DB
        ? { ...settings, mongo: payload(), type: dbType }
        : { ...settings, [dbType]: payload(), type: dbType },
    );
  };

  return (
    <div className="iap-db-connection-section">
      <form className="iap-db-form" onSubmit={handleSubmit}>
        {!!error && (
          <FormErrorMessage error={error} />
        )}

        <Field label="Database type">
          <Select
            name="type"
            value={dbType}
            disabled={posting}
            onChange={changeDbType}
            placeholder="Select Database Type"
          >
            <Option value={BOLT_DB} title="Bolt DB" />
            <Option value={MONGO_DB} title="Mongo DB" />
            <Option value={DYNAMO_DB} title="Dynamo DB" />
          </Select>
        </Field>

        {dbType === DYNAMO_DB && (
          <>
            <Field label="Region">
              <Input
                name="region"
                value={dbSettings.region}
                placeholder="e.g. ap-northeast-3"
                onChange={handleInput}
                disabled={posting}
                errorMessage={validation.region}
              />
            </Field>
            <Field label="Endpoint">
              <Input
                name="endpoint"
                value={dbSettings.endpoint}
                placeholder="e.g. localhost:27017"
                disabled={posting}
                onChange={handleInput}
                errorMessage={validation.endpoint}
              />
            </Field>
          </>
        )}

        {dbType === MONGO_DB && (
          <>
            <Field label="Name">
              <Input
                name="database"
                value={dbSettings.database}
                autoComplete="off"
                placeholder="e.g. identifo"
                disabled={posting}
                onChange={handleInput}
                errorMessage={validation.database}
              />
            </Field>
            <Field label="Endpoint">
              <Input
                name="connection"
                value={dbSettings.connection}
                autoComplete="off"
                placeholder="e.g. identifo"
                disabled={posting}
                onChange={handleInput}
                errorMessage={validation.connection}
              />
            </Field>
          </>
        )}

        {dbType === BOLT_DB && (
          <Field label="Path">
            <Input
              name="path"
              value={dbSettings.path}
              placeholder="./db.db"
              onChange={handleInput}
              disabled={posting}
              errorMessage={validation.path}
            />
          </Field>
        )}

        <footer className="iap-db-form__footer">
          <Button
            error={!posting && !!error}
            type="submit"
            Icon={posting ? LoadingIcon : SaveIcon}
            disabled={posting}
          >
            Save Changes
          </Button>
          <Button
            transparent
            disabled={posting}
            onClick={onCancel}
          >
            Cancel
          </Button>
        </footer>
      </form>
    </div>
  );
};

ConnectionSettingsForm.propTypes = {
  posting: PropTypes.bool.isRequired,
  settings: PropTypes.shape({
    type: PropTypes.string,
    mongo: PropTypes.shape({
      connection: PropTypes.string,
      database: PropTypes.string,
    }),
    boltdb: PropTypes.shape({
      path: PropTypes.string,
    }),
    dynamo: PropTypes.shape({
      region: PropTypes.string,
      endpoint: PropTypes.string,
    }),
  }),
  onCancel: PropTypes.func,
  onSubmit: PropTypes.func.isRequired,
  error: PropTypes.instanceOf(Error),
};

ConnectionSettingsForm.defaultProps = {
  settings: {
    type: '',
    mongo: {},
    boltdb: {},
    dynamo: {},
  },
  onCancel: null,
  error: null,
};

export default ConnectionSettingsForm;
