import { FC } from "react";
import { Route, Switch } from "react-router-dom";
import { Demo } from "../components";
import { useEnsureAuthentication } from "../utils/hooks/ensureAuthentication";
import App from "./App";
import Callback from './Callback';

const Root: FC = () => {
    useEnsureAuthentication()
    return (
        <Switch>
            <Route exact path='/' component={App} />
            <Route exact path='/callback' component={Callback} />
            <Route exact path='/demo' component={Demo} />
        </Switch>
    );
}

export default Root;