/* eslint-disable camelcase */
import React, { useState, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import classnames from 'classnames';
import EditUserForm from './Form';
import ActionsButton from '~/components/shared/ActionsButton';
import DropdownIcon from '~/components/icons/DropdownIcon';
import {
  fetchUserById, alterUser, deleteUserById, resetUserError,
} from '~/modules/users/actions';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const goBackPath = '/management/users';

const EditUserView = ({ match, history }) => {
  const dispatch = useDispatch();
  const id = match.params.userid;

  const user = useSelector(s => s.selectedUser.user);
  const error = useSelector(s => s.selectedUser.error);

  const [isIdsShown, setIdsShown] = useState(false);

  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const federatedIds = useMemo(() => {
    if (user && user.federated_ids) {
      return user.federated_ids;
    }
    return [];
  }, [user]);

  const fetchData = async () => {
    setProgress(70);
    await dispatch(fetchUserById(id));
    setProgress(100);
  };

  React.useEffect(() => {
    fetchData();
  }, []);

  const handleDeleteClick = async () => {
    setProgress(70);

    try {
      await dispatch(deleteUserById(id));
      notifySuccess({
        title: 'Deleted',
        text: 'User has been deleted successfully',
      });
      history.push(goBackPath);
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'User could not be updated',
      });
    } finally {
      setProgress(100);
    }
  };

  const handleSubmit = async (data) => {
    setProgress(70);

    try {
      await dispatch(alterUser(id, data));
      notifySuccess({
        title: 'Updated',
        text: 'User has been updated successfully',
      });
      history.push(goBackPath);
    } catch (_) {
      notifyFailure({
        title: 'Something went wrong',
        text: 'User could not be updated',
      });
    } finally {
      setProgress(100);
    }
  };

  const handleCancel = () => {
    dispatch(resetUserError());
    history.push(goBackPath);
  };

  const toggleDropdown = () => {
    setIdsShown(!isIdsShown);
  };

  const availableActions = [{
    title: 'Delete User',
    onClick: handleDeleteClick,
  }];

  const dropdownIconClass = classnames({
    'iap-dropdown-icon': true,
    'iap-dropdown-icon--open': isIdsShown,
  });

  const federatedKeys = classnames({
    'iap-management-section__federated-keys': true,
    'iap-management-section__federated-keys--visible': isIdsShown,
  });

  return (
    <section className="iap-management-section">
      <header>
        <div>
          <Link to={goBackPath} className="iap-management-section__back">
            ← &nbsp;Users
          </Link>
        </div>
        <div className="iap-management-section__title">
          User Details
          <ActionsButton loading={!!progress} actions={availableActions} />
        </div>
        <p className="iap-management-section__description">
          <span className="iap-section-description__id">
            id:&nbsp;
            {id}
          </span>
        </p>
        {!!federatedIds.length && (
          <div className="iap-management-section__federated">
            <p className="iap-management-section__federated-title">Federated ids</p>
            <DropdownIcon className={dropdownIconClass} onClick={toggleDropdown} />
          </div>
        )}
        <div className={federatedKeys}>
          {federatedIds.map((i) => {
            return (
              <p className="iap-management-section__description" key={i}>
                <span className="iap-section-description__id">{i}</span>
              </p>
            );
          })}
        </div>
      </header>
      <main>
        <EditUserForm
          user={user}
          error={error}
          loading={!!progress}
          onCancel={handleCancel}
          onSubmit={handleSubmit}
        />
      </main>
    </section>
  );
};

export default EditUserView;
