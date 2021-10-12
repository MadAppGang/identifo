import { FC, useContext, useEffect } from "react";
import { Redirect } from "react-router-dom";
import { Context as AppContext } from "../context/app-context";
import { identifo } from "../services/identifo";

const Callback: FC = () => {
  const { state, actions } = useContext(AppContext);
  useEffect(() => {
    (async function () {
      await identifo.handleAuthentication();
      actions.setIsAuth(identifo.isAuth);
    })();
  }, [actions]);
  if (state.isAuthenticated) {
    return <Redirect to="/demo" />;
  }
  return null;
};

export default Callback;
