import React, { useContext, useState } from "react";
import { Redirect } from "react-router-dom";
import { Header, Main } from "../components";
import { Context as AppContext } from "../context/app-context";

function App() {
  const [result, setResult] = useState<boolean>(false);
  const { state } = useContext(AppContext);
  if (state.isAuthenticated) return <Redirect to="/demo" />;
  const onComplete = () => {
    setResult(true);
  };
  if (result) return <Redirect to="/" />;
  return (
    <div className="App">
      <Header />
      <div className="form-holder">
        <identifo-form
          onComplete={onComplete}
          url="http://localhost:8081"
          app-id="c3vqvhea0brnc4dvdnvg"
        ></identifo-form>
      </div>
    </div>
  );
}

export default App;
