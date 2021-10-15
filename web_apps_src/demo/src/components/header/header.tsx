import React, { FC, useContext } from "react";
import { Context as AppContext } from "../../context/app-context";
import { identifo } from "../../services/identifo";

const Header: FC = React.memo(() => {
  const { state } = useContext(AppContext);
  const login = () => {
    identifo.signin();
  };
  const signUp = () => {
    identifo.signup();
  };
  const logout = () => {
    identifo.logout();
  };

  return (
    <header className="header">
      <div className="container">
        <div className="header_in">
          <div className="header_logo">
            <span>identifo.js</span>
            <span>Example login with web flow -&gt; </span>
          </div>
          <div className="header_buttons">
            {!state.isAuthenticated ? (
              <>
                <button className="outlined-btn" onClick={login}>
                  Login
                </button>
                <button className="primary-btn" onClick={signUp}>
                  Sign up
                </button>
              </>
            ) : (
              <>
                <button className="primary-btn" onClick={logout}>
                  Logout
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
});

export default Header;
