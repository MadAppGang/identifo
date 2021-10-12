import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { login, checkAuthState } from '~/modules/auth/actions';
import Input from '~/components/shared/Input';
import Button from '~/components/shared/Button';
import LoadingIcon from '~/components/icons/LoadingIcon';
import EmailIcon from '~/components/icons/EmailIcon';
import PasswordIcon from '~/components/icons/PasswordIcon';
import FormErrorMessage from '~/components/shared/FormErrorMessage';
import useForm from '~/hooks/useForm';

const LoginForm = () => {
  const dispatch = useDispatch();
  const signingIn = useSelector(s => s.auth.inProgress);
  const error = useSelector(s => s.auth.error);

  const handleSubmit = (values) => {
    dispatch(login(values.email, values.password));
  };

  React.useEffect(() => {
    dispatch(checkAuthState());
  }, []);

  const form = useForm({
    email: '',
    password: '',
  }, null, handleSubmit);

  return (
    <form className="iap-login-form" onSubmit={form.handleSubmit}>
      <h1 className="login-form__logo">
        <span>identifo</span>
        <span>Admin Panel</span>
      </h1>

      {error && (
        <div className="iap-login-form__err">
          <FormErrorMessage error={error} />
        </div>
      )}

      <Input
        name="email"
        value={form.values.email}
        placeholder="Email"
        disabled={signingIn}
        Icon={EmailIcon}
        onChange={form.handleChange}
      />

      <Input
        name="password"
        type="password"
        value={form.values.password}
        placeholder="Password"
        disabled={signingIn}
        Icon={PasswordIcon}
        onChange={form.handleChange}
      />

      <footer className="iap-login-form__footer">
        <Button
          stretch
          type="submit"
          disabled={signingIn}
          Icon={signingIn ? LoadingIcon : null}
          error={!!error}
        >
          Sign In
        </Button>
      </footer>
    </form>
  );
};

LoginForm.defaultProps = {
  signingIn: false,
  error: null,
};

export default LoginForm;
