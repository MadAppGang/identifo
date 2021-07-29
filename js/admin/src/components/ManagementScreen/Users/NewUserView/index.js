import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import UserForm from './UserForm';
import { postUser, resetUserError } from '~/modules/users/actions';
import useProgressBar from '~/hooks/useProgressBar';
import useNotifications from '~/hooks/useNotifications';

const goBackPath = '/management/users';

const NewUserView = ({ history }) => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();
  const { notifySuccess, notifyFailure } = useNotifications();

  const user = useSelector(s => s.selectedUser.user);
  const error = useSelector(s => s.selectedUser.error);

  React.useEffect(() => {
    if (user && user.id && progress === 100) {
      history.push(`/management/users/${user.id}`);
    }
  }, [user, progress]);

  const handleSubmit = async (data) => {
    setProgress(70);

    try {
      await dispatch(postUser(data));

      notifySuccess({
        title: 'Created',
        text: 'User has been created successfully',
      });
    } catch (_) {
      notifyFailure({
        title: 'Error',
        text: 'User could not be created',
      });
    } finally {
      setProgress(100);
    }
  };

  const handleCancel = () => {
    dispatch(resetUserError());
    history.push(goBackPath);
  };

  return (
    <section className="iap-management-section">
      <header>
        <div>
          <Link to={goBackPath} className="iap-management-section__back">
            ← &nbsp;Users
          </Link>
        </div>
        <p className="iap-management-section__title">
          Create User
        </p>
        <p className="iap-management-section__description">
          Created user is going to be able to log in using these credentials
        </p>
      </header>
      <main>
        <UserForm
          error={error}
          saving={!!progress}
          onCancel={handleCancel}
          onSubmit={handleSubmit}
        />
      </main>
    </section>
  );
};

export default NewUserView;
