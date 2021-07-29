import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import UsersPlaceholder from './Placeholder';
import UserList from './UserList';
import { fetchUsers } from '~/modules/users/actions';
import Button from '~/components/shared/Button';
import AddIcon from '~/components/icons/AddIcon';
import SearchInput from '~/components/shared/SearchInput';

class UsersMainView extends Component {
  constructor() {
    super();

    this.state = {
      searchQuery: '',
    };

    this.handleSearch = this.handleSearch.bind(this);
  }

  componentDidMount() {
    this.props.fetchUsers();
  }

  handleSearch(searchQuery) {
    this.props.fetchUsers({ search: searchQuery });
    this.setState({ searchQuery });
  }

  render() {
    const { users, fetching, history } = this.props;
    const { searchQuery } = this.state;

    if (!users.length && !fetching && !searchQuery) {
      return (
        <section className="iap-management-section">
          <UsersPlaceholder
            onCreateNewUserClick={() => history.push('/management/users/new')}
          />
        </section>
      );
    }

    return (
      <section className="iap-management-section">
        <p className="iap-management-section__title">
          Users
          <Button
            Icon={AddIcon}
            onClick={() => history.push('/management/users/new')}
          >
            Create User
          </Button>
        </p>

        <p className="iap-management-section__description">
          Look for users, edit, delete them and add new ones.
        </p>

        <SearchInput
          timeout={400}
          placeholder="Search for users"
          onChange={this.handleSearch}
        />

        <UserList loading={fetching} users={users} />
      </section>
    );
  }
}

UsersMainView.propTypes = {
  fetchUsers: PropTypes.func.isRequired,
  users: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.string,
    name: PropTypes.string,
    email: PropTypes.string,
    latest_login_time: PropTypes.number,
    num_of_logins: PropTypes.number,
  })),
  fetching: PropTypes.bool,
  history: PropTypes.shape({
    push: PropTypes.func,
  }).isRequired,
};

UsersMainView.defaultProps = {
  users: [],
  fetching: false,
};

const mapStateToProps = state => ({
  fetching: state.users.fetching,
  users: state.users.list,
});

const actions = {
  fetchUsers,
};

export { UsersMainView };

export default connect(mapStateToProps, actions)(UsersMainView);
