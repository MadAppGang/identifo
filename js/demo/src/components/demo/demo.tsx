import { FC, useContext, useState } from "react";
import { Redirect } from "react-router-dom";
import { Header } from "../index";
import { Context as AppContext } from "../../context/app-context";
import { identifo } from "../../services/identifo";

const Demo: FC = () => {
  const [result, setResult] = useState<any>(null);
  const { state } = useContext(AppContext);
  const renewSession = () => {
    identifo
      .renewSession()
      .then((t) => console.log(t))
      .catch((e) => console.log(e));
  };

  const getToken = async () => {
    const token = await identifo.getToken();
    setResult(JSON.stringify(token, undefined, 4));
  };
  const getAuth = () => {
    setResult(JSON.stringify(identifo.isAuth, undefined, 4));
  };

  if (!state.isAuthenticated) return <Redirect to="/" />;
  return (
    <>
      <Header />
      <div className="demo">
        <div className="demo_buttons">
          <button className="primary-btn" onClick={getAuth}>
            Get Authenticated
          </button>
          <button className="primary-btn" onClick={getToken}>
            Get Token
          </button>
          <button className="primary-btn" onClick={renewSession}>
            Renew Session
          </button>
        </div>
        <div className="demo_result">{result && <pre>{result}</pre>}</div>
      </div>
    </>
  );
};

export default Demo;
