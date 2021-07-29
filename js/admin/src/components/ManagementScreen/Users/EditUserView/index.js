import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import EditUserForm from './Form';
import ActionsButton from '~/components/shared/ActionsButton';
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

  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

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

  const availableActions = [{
    title: 'Delete User',
    onClick: handleDeleteClick,
  }];

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
