import { useContext, useEffect } from "react";
import { Context as AppContext } from "../../context/app-context";
import { identifo } from "../../services/identifo";

export const useEnsureAuthentication = () => {
  const { actions } = useContext(AppContext);

  useEffect(() => {
    const status = identifo.isAuth;
    actions.setIsAuth(status);
  }, [actions]);
};
