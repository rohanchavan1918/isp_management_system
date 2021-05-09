import React from "react"
import {
    BrowserRouter as Router,
    Switch,
    Route,
  } from "react-router-dom";

import Signin from '../pages/auth/Signin'
import Plans from '../pages/Plans'

const ispRouter = () => {
    return (
        <div>
        <Router>
            <Switch>
                <Route exact path="/" >
                    <Signin />
                </Route>
                <Route exact path="/signin" >
                    <Signin />
                </Route>
                <Route exact path="/plans" >
                    <Plans />
                </Route>
            </Switch>
        </Router>
        </div>
    )
}

export default ispRouter