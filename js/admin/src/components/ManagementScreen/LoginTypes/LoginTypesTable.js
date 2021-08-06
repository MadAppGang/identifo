import React from 'react';
import { DatagridRow, DatagridHeader } from '~/components/shared/Datagrid';
import Toggle from '~/components/shared/Toggle';
import PhoneLoginIcon from '~/components/icons/PhoneLoginIcon.svg';
import FederatedLoginIcon from '~/components/icons/FederatedLoginIcon.svg';
import UsernameLoginIcon from '~/components/icons/UsernameLoginIcon.svg';

const datagrid = {
  icon: {
    title: '',
    width: '13%',
  },
  type: {
    title: 'Login With',
    width: '62%',
  },
  enabled: {
    title: 'Enabled',
    width: '25%',
  },
};

const LoginTypesTable = (props) => {
  const { types, onChange } = props;

  return (
    <div className="login-types-table">
      <DatagridHeader>
        {Object.keys(datagrid).map(key => (
          <div key={key} style={{ width: datagrid[key].width }}>
            {datagrid[key].title}
          </div>
        ))}
      </DatagridHeader>

      <DatagridRow className="login-types-row">
        <div style={{ width: datagrid.icon.width }}>
          <div className="login-types-row__icon">
            <UsernameLoginIcon className="login-type-icon" />
          </div>
        </div>
        <div style={{ width: datagrid.type.width }}>
          <p className="login-types-row__type">Username and Password</p>
        </div>
        <div style={{ width: datagrid.enabled.width }}>
          <div className="login-types-row__enabled">
            <Toggle value={types.username} onChange={v => onChange('username', v)} />
          </div>
        </div>
      </DatagridRow>

      <DatagridRow className="login-types-row">
        <div style={{ width: datagrid.icon.width }}>
          <div className="login-types-row__icon">
            <UsernameLoginIcon className="login-type-icon" />
          </div>
        </div>
        <div style={{ width: datagrid.type.width }}>
          <p className="login-types-row__type">Email and Password</p>
        </div>
        <div style={{ width: datagrid.enabled.width }}>
          <div className="login-types-row__enabled">
            <Toggle value={types.email} onChange={v => onChange('email', v)} />
          </div>
        </div>
      </DatagridRow>

      <DatagridRow className="login-types-row">
        <div style={{ width: datagrid.icon.width }}>
          <div className="login-types-row__icon">
            <PhoneLoginIcon className="login-type-icon" />
          </div>
        </div>
        <div style={{ width: datagrid.type.width }}>
          <p className="login-types-row__type">Phone Number with password or OTP</p>
        </div>
        <div style={{ width: datagrid.width }}>
          <div className="login-types-row__enabled">
            <Toggle value={types.phone} onChange={v => onChange('phone', v)} />
          </div>
        </div>
      </DatagridRow>

      <DatagridRow className="login-types-row">
        <div style={{ width: datagrid.icon.width }}>
          <div className="login-types-row__icon">
            <FederatedLoginIcon className="login-type-icon" />
          </div>
        </div>
        <div style={{ width: datagrid.type.width }}>
          <p className="login-types-row__type">Federated Identity</p>
        </div>
        <div style={{ width: datagrid.enabled.width }}>
          <div className="login-types-row__enabled">
            <Toggle value={types.federated} onChange={v => onChange('federated', v)} />
          </div>
        </div>
      </DatagridRow>
    </div>
  );
};

export default LoginTypesTable;
