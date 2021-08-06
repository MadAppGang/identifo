import React, { Component } from 'react';
import PropTypes from 'prop-types';
import update from '@madappgang/update-by-path';
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

class ConnectionSettingsForm extends Component {
  constructor({ settings }) {
    super();

    this.validate = Validation.applyRules(databaseFormValidationRules);

    this.state = {
      settings,
      validation: {
        type: '',
        endpoint: '',
        name: '',
        region: '',
        path: '',
      },
    };

    this.handleInput = this.handleInput.bind(this);
    this.handleDBTypeChange = this.handleDBTypeChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleBlur = this.handleBlur.bind(this);
  }

  getFieldsToOmitDuringValidation() {
    switch (this.state.settings.type) {
      case DYNAMO_DB: return ['name', 'path', 'endpoint'];
      case MONGO_DB: return ['region', 'path'];
      case BOLT_DB: return ['name', 'region', 'endpoint'];
      default: return [];
    }
  }

  handleInput({ target }) {
    const { name, value } = target;
    let { validation } = this.state;

    if (validation[name]) {
      validation = update(validation, { [name]: '' });
    }

    this.setState(state => ({
      settings: update(state.settings, {
        [name]: value,
      }),
      validation,
    }));
  }

  handleBlur({ target }) {
    const { name, value } = target;
    const validationMessage = this.validate(name, {
      ...this.state.settings,
      [name]: value,
    });

    this.setState(state => ({
      validation: update(state.validation, {
        [name]: validationMessage,
      }),
    }));
  }

  handleDBTypeChange(type) {
    this.setState(state => ({
      settings: update(state.settings, { type }),
      validation: Validation.reset(state.validation),
    }));
  }

  handleSubmit(event) {
    event.preventDefault();

    const validation = this.validate('all', this.state.settings, {
      omit: this.getFieldsToOmitDuringValidation(),
    });

    if (Validation.hasError(validation)) {
      this.setState({ validation });
      return;
    }

    let { settings } = this.state;

    settings = update(settings, {
      region: region => settings.type === DYNAMO_DB ? region : '',
      name: name => settings.type === MONGO_DB ? name : '',
      path: path => settings.type === BOLT_DB ? path : '',
      endpoint: endpoint => settings.type !== BOLT_DB ? endpoint : '',
    });

    this.props.onSubmit(settings);
  }

  render() {
    const { settings, validation } = this.state;
    const { posting, error } = this.props;
    const { type, name, region, endpoint, path } = settings;

    return (
      <div className="iap-db-connection-section">
        <form className="iap-db-form" onSubmit={this.handleSubmit}>
          {!!error && (
            <FormErrorMessage error={error} />
          )}

          <Field label="Database type">
            <Select
              name="type"
              value={type}
              disabled={posting}
              onChange={this.handleDBTypeChange}
              placeholder="Select Database Type"
            >
              <Option value={BOLT_DB} title="Bolt DB" />
              <Option value={MONGO_DB} title="Mongo DB" />
              <Option value={DYNAMO_DB} title="Dynamo DB" />
            </Select>
          </Field>

          {type === DYNAMO_DB && (
            <Field label="Region">
              <Input
                name="region"
                value={region}
                placeholder="e.g. ap-northeast-3"
                onChange={this.handleInput}
                disabled={posting}
                errorMessage={validation.region}
                onBlur={this.handleBlur}
              />
            </Field>
          )}

          {type === MONGO_DB && (
            <Field label="Name">
              <Input
                name="name"
                value={name}
                autoComplete="off"
                placeholder="e.g. identifo"
                disabled={posting}
                onChange={this.handleInput}
                errorMessage={validation.name}
                onBlur={this.handleBlur}
              />
            </Field>
          )}

          {type === BOLT_DB && (
            <Field label="Path">
              <Input
                name="path"
                value={path}
                placeholder="./db.db"
                onChange={this.handleInput}
                disabled={posting}
                errorMessage={validation.path}
                onBlur={this.handleBlur}
              />
            </Field>
          )}

          {type !== BOLT_DB && (
            <Field label="Endpoint">
              <Input
                name="endpoint"
                value={endpoint}
                placeholder="e.g. localhost:27017"
                disabled={posting}
                onChange={this.handleInput}
                onBlur={this.handleBlur}
                errorMessage={validation.endpoint}
              />
            </Field>
          )}

          <footer className="iap-db-form__footer">
            <Button
              error={!posting && !!error}
              type="submit"
              Icon={posting ? LoadingIcon : SaveIcon}
              disabled={posting || Validation.hasError(validation)}
            >
              Save Changes
            </Button>
            <Button
              transparent
              disabled={posting}
              onClick={this.props.onCancel}
            >
              Cancel
            </Button>
          </footer>
        </form>
      </div>
    );
  }
}

ConnectionSettingsForm.propTypes = {
  posting: PropTypes.bool.isRequired,
  settings: PropTypes.shape({
    type: PropTypes.string,
    endpoint: PropTypes.string,
    name: PropTypes.string,
    region: PropTypes.string,
  }),
  onCancel: PropTypes.func,
  onSubmit: PropTypes.func.isRequired,
  error: PropTypes.instanceOf(Error),
};

ConnectionSettingsForm.defaultProps = {
  settings: {
    type: '',
    endpoint: '',
    name: '',
    region: '',
  },
  onCancel: null,
  error: null,
};

export default ConnectionSettingsForm;
