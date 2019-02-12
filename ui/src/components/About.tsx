import * as React from "react";
import { RouteComponentProps } from "react-router";
import { Link, Route, Switch } from "react-router-dom";

const AboutTop = () => (
    <div>
        <h2>About Top</h2>
    </div>
);

const AboutMe = () => (
    <div>
        <h2>About Me</h2>
    </div>
);

const AboutUs = () => (
    <div>
        <h2>About Us</h2>
    </div>
);

const AboutNothing = () => (
    <div>
        <h2>About Nothing</h2>
    </div>
);

export const About = (p: RouteComponentProps) => (
    <div>
        <h2>About</h2>
        <ul>
            <li><Link to={`${p.match.url}/`}>Top</Link></li>
            <li><Link to={`${p.match.url}/me`}>Me</Link></li>
            <li><Link to={`${p.match.url}/us`}>Us</Link></li>
            <li><Link to={`${p.match.url}/abc`}>ABC</Link></li>
        </ul>
        <Switch>
            <Route component={AboutTop} path={`${p.match.url}/`} exact={true} strict={true} />
            <Route component={AboutMe} path={`${p.match.url}/me`} exact={true} strict={true} />
            <Route component={AboutUs} path={`${p.match.url}/us`} exact={true} strict={true} />
            <Route component={AboutNothing} />
        </Switch>
    </div>
);
