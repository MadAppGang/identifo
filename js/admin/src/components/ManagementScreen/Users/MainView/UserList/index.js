import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import UserRow from './UserRow';
import UserHeader from './UserHeader';
import Preloader from './Preloader';
import { DatagridNotFound } from '~/components/shared/Datagrid';

import './index.css';

const datagrid = {
  icon: {
    title: '',
    width: '10%',
  },
  name: {
    title: 'Username',
    width: '30%',
  },
  email: {
    title: 'Email',
    width: '30%',
  },
  phone: {
    title: 'Phone',
    width: '25%',
  },
  num_of_logins: {
    title: '# of Logins',
    width: '10%',
  },
};

const UserList = (props) => {
  const { users, loading } = props;

  return (
    <div className="iap-userlist">
      <div className="iap-datagrid">
        <UserHeader config={datagrid} />
        <main className="iap-datagrid-body">
          {loading && (
            <Preloader />
          )}

          {!loading && users.map(user => (
            <Link key={user.id} to={`/management/users/${user.id}`} className="rrdl">
              <UserRow data={user} config={datagrid} />
            </Link>
          ))}

          {!users.length && !loading && (
            <DatagridNotFound text="No Users Found" />
          )}
        </main>
      </div>
    </div>
  );
};

UserList.propTypes = {
  users: PropTypes.arrayOf(PropTypes.shape()),
  loading: PropTypes.bool,
};

UserList.defaultProps = {
  users: [],
  loading: false,
};

export default UserList;
